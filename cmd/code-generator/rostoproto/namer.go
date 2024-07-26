package rostoproto

import (
	"fmt"
	"strings"
)

type Namer interface {
	Name(*Type) string
}

type Names map[*Type]string

type NameSystems map[string]Namer

type NameStrategy struct {
	Prefix, Suffix string
	Join           func(pre string, parts []string, post string) string
	IgnoreWords    map[string]bool
	// For example, if Ignore words lists "proto" and type Foo is in
	// pkg/server/frobbing/proto, then a value of 1 will give a type name
	// of FrobbingFoo, 2 gives ServerFrobbingFoo, etc.
	PrependPackageNames int
	// a cache of names this Namer has assigned
	Names
}

type ImportTracker interface {
	AddType(*Type)
	AddSymbol(Name)
	LocalNameOf(packagePath string) string
	PathOf(localName string) (string, bool)
	ImportLines() []string
}

// helper function which returns a namer that outputs camelcase names
func NewNamer(prependPackageNames int, ignoreWords ...string) *NameStrategy {
	n := &NameStrategy{
		Join:                Joiner(IC, IC),
		IgnoreWords:         map[string]bool{},
		PrependPackageNames: prependPackageNames,
	}

	for _, w := range ignoreWords {
		n.IgnoreWords[w] = true
	}
	return n
}

func IC(in string) string {
	if in == "" {
		return in
	}
	return strings.ToUpper(in[:1]) + in[1:]
}

func IL(in string) string {
	if in == "" {
		return in
	}
	return strings.ToLower(in[:1]) + in[1:]
}

// specify functions which preprocess various components of a name before joining them. eg: IC & IL
func Joiner(first, others func(string) string) func(pre string, in []string, post string) string {
	return func(pre string, in []string, post string) string {
		tmp := []string{others(pre)}
		for i := range in {
			tmp = append(tmp, others(in[i]))
		}
		tmp = append(tmp, others(post))
		return first(strings.Join(tmp, ""))
	}
}

func (ns *NameStrategy) removePrefixAndSuffix(s string) string {
	lowerIn := strings.ToLower(s)
	lowerP := strings.ToLower(ns.Prefix)
	lowerS := strings.ToLower(ns.Suffix)
	p, e := 0, len(s)
	if strings.HasPrefix(lowerIn, lowerP) {
		p = len(ns.Prefix)
	}
	if strings.HasSuffix(lowerIn, lowerS) {
		e -= len(ns.Suffix)
	}
	return s[p:e]

}

func (ns *NameStrategy) Name(t *Type) string {
	if ns.Names == nil {
		ns.Names = Names{}
	}

	if s, ok := ns.Names[t]; ok {
		return s
	}

	// lets not add the name of the package in the name for now
	name := ns.Join(ns.Prefix, []string{t.Name.Name}, ns.Suffix)
	ns.Names[t] = name
	return name
}

func protoSafePackage(name string) string {
	pkg := strings.ReplaceAll(name, "/", ".")
	return strings.ReplaceAll(pkg, "-", "_")
}

func protoSafePackageMsg(name string) string {
	msg := strings.Split(name, ".")[0]
	return protoSafePackage(msg)
}

type localNamer struct {
	localPackage Name
}

func (n localNamer) Name(t *Type) string {
	if len(n.localPackage.Package) != 0 && n.localPackage.Package == t.Name.Package {
		return t.Name.Name
	}

	if strings.Contains(t.Name.Package, ".") {
		return fmt.Sprintf(".%s", t.Name)
	}

	return t.Name.String()
}

type ProtobufNamer struct {
	packages []*ProtobufPackage
	// the key for this will be a ros msg import path
	packagesByPath map[string]*ProtobufPackage
}

func NewProtobufNamer() *ProtobufNamer {
	return &ProtobufNamer{
		packagesByPath: make(map[string]*ProtobufPackage),
	}
}

func (n *ProtobufNamer) Name(t *Type) string {
	return t.Name.String()
}

func (n *ProtobufNamer) Add(p *ProtobufPackage) {
	if _, ok := n.packagesByPath[p.Path()]; !ok {
		n.packagesByPath[p.Path()] = p
		n.packages = append(n.packages, p)
	}
}

func (n *ProtobufNamer) RosNameToProtoName(name Name) Name {
	if p, ok := n.packagesByPath[name.Package]; ok {
		return Name{
			Name:    name.Name,
			Package: p.Name(),
			Path:    p.ImportPath(),
		}
	}

	return Name{Name: name.Name}
}

type typeNameset map[Name]*ProtobufPackage
