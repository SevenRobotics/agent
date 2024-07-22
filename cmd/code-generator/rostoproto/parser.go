package rostoproto

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Parser struct {

	//Map of package paths to definitions.
	RosPkgs map[string]*Package
	//keep track of which packages were directly requested
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

// load packages should parse and hold the message definitions of the loaded packages?
func (p *Parser) LoadPackages(patterns ...string) ([]*Package, error) {
	existingPkgs, netNewPkgs, err := p.alreadyLoaded(patterns...)
	if err != nil {
		return nil, err
	}

	if len(netNewPkgs) == 0 {
		return existingPkgs, nil
	}

	for _, pkgPath := range netNewPkgs {
		pkg, err := p.load(pkgPath)
		if err != nil {
			continue
		}
		if pkg != nil {
			existingPkgs = append(existingPkgs, pkg)
		}
	}

	return existingPkgs, nil
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
			pkg.MessageDefs = append(pkg.MessageDefs, msgDef)
		}
		return nil
	})

	return &pkg, nil
}
