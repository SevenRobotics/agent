package rostoproto

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Parser struct {

	//Map of package paths to definitions.
	RosPkgs map[string]*Package
	//keep track of which packages were directly requested
	//used later on while loading to the generator universe
	userRequested map[string]bool

	fullyProcessed map[string]bool

	buildTags []string
}

type Options struct {
	BuildTags []string
}

func NewParser() *Parser {
	return NewWithOptions(Options{})
}

func NewWithOptions(opts Options) *Parser {
	return &Parser{
		RosPkgs:        map[string]*Package{},
		userRequested:  map[string]bool{},
		fullyProcessed: map[string]bool{},
		buildTags:      opts.BuildTags,
	}
}

func (p *Parser) alreadyLoaded(patterns ...string) ([]*Package, []string, error) {
	existingPkgs := make([]*Package, 0, len(patterns))
	netNewPkgs := make([]string, 0, len(patterns))

	// Expand and canonicalize the requested patterns.  This should be fast.
	for _, pkgPath := range patterns {
		if pkg := p.RosPkgs[pkgPath]; pkg != nil {
			existingPkgs = append(existingPkgs, pkg)
		} else {
			netNewPkgs = append(netNewPkgs, pkgPath)
		}
	}
	return existingPkgs, netNewPkgs, nil
}

func (p *Parser) LoadPackages(patterns ...string) error {
	_, err := p.loadPackages(patterns...)
	return err
}

func (p *Parser) LoadPackagesTo(u *Universe, patterns ...string) ([]*Package, error) {
	pkgs, err := p.loadPackages(patterns...)

	if err != nil {
		return nil, err
	}

	if err := p.addPkgsToUniverse(pkgs, u); err != nil {
		return nil, err
	}

	return pkgs, nil
}

// load packages should parse and hold the message definitions of the loaded packages?
func (p *Parser) loadPackages(patterns ...string) ([]*Package, error) {
	log.Printf("loading packages: %v", patterns)
	existingPkgs, netNewPkgs, err := p.alreadyLoaded(patterns...)

	if err != nil {
		return nil, err
	}

	if len(existingPkgs) > 0 {
		keys := make([]string, 0, len(existingPkgs))
		for _, p := range existingPkgs {
			keys = append(keys, p.Dir)
		}
		log.Printf("Already have: %v", existingPkgs)
	}

	for _, pkg := range existingPkgs {
		if !p.userRequested[pkg.Dir] {
			p.userRequested[pkg.Dir] = true
		}
	}

	if len(netNewPkgs) == 0 {
		return existingPkgs, nil
	}

	log.Printf("To be loaded %v", netNewPkgs)

	for _, pkg := range netNewPkgs {
		if !p.userRequested[pkg] {
			p.userRequested[pkg] = true
		}
	}

	newPkgs := []*Package{}
	for _, pkgPath := range netNewPkgs {
		pkg, err := p.load(pkgPath)
		if err != nil {
			continue
		}
		if pkg != nil {
			newPkgs = append(newPkgs, pkg)
		}
	}
	absorbPkgs := func(pkg *Package) error {
		p.RosPkgs[pkg.Dir] = pkg
		// have to take care of imports and stuff
		// hence following the error returning signature from gengo v2
		return nil
	}

	if err := forEachPackageRecursive(newPkgs, absorbPkgs); err != nil {
		return nil, err
	}
	return append(existingPkgs, newPkgs...), nil
}

func (p *Parser) load(pkgPath string) (*Package, error) {
	// get ros package name from the pkg path
	// check if this package has messages
	// walk through the folder and read every message type
	// for every message type parse it and append to Package.MessageDefinitions
	// If any error collect errors or return ?
	pkgPath = filepath.Clean(pkgPath)
	pkg := Package{}
	if _, err := os.Stat(pkgPath + "/msg"); os.IsNotExist(err) {
		return nil, err
	}
	tmp := strings.Split(pkgPath, "/")
	pkgname := tmp[len(tmp)-1]
	pkg.Dir = pkgPath
	pkg.Name = pkgname
	pkg.Path = pkgname // ros uses names as imports
	pkg.MessageDefs = map[string]map[string]*MessageDefinition{}

	filepath.WalkDir(pkgPath+"/msg", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ".msg" {
			name := strings.TrimSuffix(filepath.Base(path), ".msg")
			buf, err := os.ReadFile(path)

			if err != nil {
				return err
			}

			content := string(buf)
			msgDef, err := ParseMessageDefinition(pkgname, name, content)
			if err != nil {
				return err
			}
			if _, ok := pkg.MessageDefs[pkgname]; !ok {
				pkg.MessageDefs[pkgname] = map[string]*MessageDefinition{}
			}
			t := pkg.MessageDefs[pkgname]
			t[name] = msgDef
			pkg.MessageDefs[pkgname] = t
		}
		return nil
	})
	// log.Printf("Parsed %v messages in %v package: %v", len(pkg.MessageDefs[pkgname]),
	// 	pkgname, pkg.MessageDefs[pkgname])
	return &pkg, nil
}

func forEachPackageRecursive(pkgs []*Package, fn func(pkg *Package) error) error {
	seen := map[string]bool{}
	var errs []error
	for _, pkg := range pkgs {
		errs = append(errs, recursePackage(pkg, fn, seen)...)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func recursePackage(pkg *Package, fn func(pkg *Package) error, seen map[string]bool) []error {
	if pkg == nil {
		return nil
	}
	if seen[pkg.Dir] {
		return nil
	}
	var errs []error
	seen[pkg.Dir] = true
	if err := fn(pkg); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (p *Parser) userRequestedPackages() []string {
	pkgPaths := make([]string, 0, len(p.userRequested))
	for k := range p.userRequested {
		pkgPaths = append(pkgPaths, string(k))
	}

	sort.Strings(pkgPaths)
	return pkgPaths
}

func (p *Parser) NewUniverse() (Universe, error) {
	u := Universe{}
	pkgs := []*Package{}

	for _, path := range p.userRequestedPackages() {
		pkgs = append(pkgs, p.RosPkgs[path])
	}

	if err := p.addPkgsToUniverse(pkgs, &u); err != nil {
		return nil, err
	}

	return u, nil
}

func (p *Parser) addPkgsToUniverse(pkgs []*Package, u *Universe) error {
	addOne := func(pkg *Package) error {
		if err := p.addPkgToUniverse(pkg, u); err != nil {
			return err
		}
		return nil
	}

	if err := forEachPackageRecursive(pkgs, addOne); err != nil {
		return err
	}
	return nil
}

func (p *Parser) addPkgToUniverse(pkg *Package, u *Universe) error {
	newU := *u
	pkgPath := pkg.Dir
	if p.fullyProcessed[pkgPath] {
		return nil
	}
	pkgPath = filepath.Base(pkgPath)

	// log.Printf("Adding package with path %v to universe", pkgPath)
	p.fullyProcessed[pkgPath] = true
	newU[pkgPath] = pkg
	*u = newU
	return nil
}
