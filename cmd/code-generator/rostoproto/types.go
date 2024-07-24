package rostoproto

import "strings"

// ref: gengo v2 types
type Name struct {
	// Empty if embedded or builtin. This is the package path unless Path is specified.
	Package string
	// The type name.
	Name string
	// An optional location of the type definition for languages that can have disjoint
	// packages and paths.
	Path string
}

// return Fully qualified Name in string format
func (N Name) String() string {
	if N.Package == "" {
		return N.Name
	}
	return N.Package + "." + N.Name
}

// return Name for a fully qualified name string, eg: uavcan_ros_bridge.Power.Battery returns
// Name {Package: uavcan_ros_bridge.Power, name: Battery}
func ParseFullyQualifiedName(n string) Name {
	cs := strings.Split(n, ".")
	pkg := ""
	if len(cs) > 1 {
		pkg = strings.Join(cs[0:len(cs)-1], ".")
	}
	return Name{
		Name:    cs[len(cs)-1],
		Package: pkg,
	}
}

const (
	Builtin    string = "Builtin"
	Array      string = "Array"
	TypeParam  string = "TypeParam"
	Protobuf   string = "Protobuf"
	RosMsgType string = "RosMsg"
)

type Package struct {
	// canonical import path of this package. eg: for proto : "google/protobuf/Empty.proto"
	Path string

	// the location of this package on the disk.
	Dir string

	//Short name of this package eg: sensor_msgs
	Name string

	MessageDefs map[string]map[string]*MessageDefinition // can add other types of defs as we go along
	//Types within this package, indexed by their name. Not sure if we need it, leaving it here
	// since it might be needed later on while generating subscribers & publishers
	Types map[string]*Type

	//Packages imported by this package. key value being the canonical package path
	Imports map[string]*Package
}

func (p *Package) HasImport(packageName string) bool {
	_, has := p.Imports[packageName]
	return has
}

func (p *Package) Type(typeName string) *Type {
	if t, ok := p.Types[typeName]; ok {
		return t
	}

	if p.Path == "" {
		if t, ok := Builtins.Types[typeName]; ok {
			p.Types[typeName] = t
			return t
		}
	}

	t := &Type{Name: Name{Name: typeName, Package: p.Path}}
	p.Types[typeName] = t
	return t
}

// universe is a map of all packages indexed by the package name
type Universe map[string]*Package

// Package returns Package for the given path
func (u Universe) Package(packagePath string) *Package {
	if p, ok := u[packagePath]; ok {
		return p
	}

	p := &Package{
		Path:    packagePath,
		Imports: map[string]*Package{},
	}
	u[packagePath] = p
	return p
}

// Register import types for a given package
func (u Universe) AddImports(packagePath string, importPaths ...string) {
	p := u.Package(packagePath)
	for _, i := range importPaths {
		p.Imports[i] = u.Package(i)
	}
}

func (u Universe) Type(typeName Name) *Type {
	return u.Package(typeName.Name).Type(typeName.Name)
}

type Type struct {
	Name Name
	Kind string // builtin primitive or array
}

func (t *Type) String() string {
	if t == nil {
		return ""
	}
	return t.Name.String()
}

func (t *Type) isBuiltin() bool {
	if t.Kind == Builtin {
		return true
	}
	return false
}

// Built in types.
var (
	String = &Type{
		Name: Name{Name: "string"},
		Kind: Builtin,
	}
	Int64 = &Type{
		Name: Name{Name: "int64"},
		Kind: Builtin,
	}
	Int32 = &Type{
		Name: Name{Name: "int32"},
		Kind: Builtin,
	}
	Int16 = &Type{
		Name: Name{Name: "int16"},
		Kind: Builtin,
	}
	Int = &Type{
		Name: Name{Name: "int"},
		Kind: Builtin,
	}
	Uint64 = &Type{
		Name: Name{Name: "uint64"},
		Kind: Builtin,
	}
	Uint32 = &Type{
		Name: Name{Name: "uint32"},
		Kind: Builtin,
	}
	Uint16 = &Type{
		Name: Name{Name: "uint16"},
		Kind: Builtin,
	}
	Uint = &Type{
		Name: Name{Name: "uint"},
		Kind: Builtin,
	}
	Uintptr = &Type{
		Name: Name{Name: "uintptr"},
		Kind: Builtin,
	}
	Float64 = &Type{
		Name: Name{Name: "float64"},
		Kind: Builtin,
	}
	Float32 = &Type{
		Name: Name{Name: "float32"},
		Kind: Builtin,
	}
	Float = &Type{
		Name: Name{Name: "float"},
		Kind: Builtin,
	}
	Bool = &Type{
		Name: Name{Name: "bool"},
		Kind: Builtin,
	}
	Byte = &Type{
		Name: Name{Name: "byte"},
		Kind: Builtin,
	}

	Builtins = &Package{
		Types: map[string]*Type{
			"bool":    Bool,
			"string":  String,
			"int":     Int,
			"int64":   Int64,
			"int32":   Int32,
			"int16":   Int16,
			"int8":    Byte,
			"uint":    Uint,
			"uint64":  Uint64,
			"uint32":  Uint32,
			"uint16":  Uint16,
			"uint8":   Byte,
			"uintptr": Uintptr,
			"byte":    Byte,
			"float":   Float,
			"float64": Float64,
			"float32": Float32,
		},
		Imports: map[string]*Package{},
		Path:    "",
		Name:    "",
	}
)

func IsInteger(t *Type) bool {
	switch t {
	case Int, Int64, Int32, Int16, Uint, Uint64, Uint32, Uint16, Byte:
		return true
	default:
		return false
	}
}
