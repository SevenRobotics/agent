package rostoproto

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"slices"
	"strings"
)

var (
	errUnrecognizedType = fmt.Errorf("did not recognize the provided type")
)

type genProtoIDL struct {
	RosMsgGenerator
	localPackage    Name
	localRosPackage Name
	imports         ImportTracker

	generateAll    bool
	omitFieldTypes map[string]struct{}
}

func (g *genProtoIDL) PackageVars(c *Context) []string {
	return []string{
		fmt.Sprintf("package %s;\n", g.localRosPackage),
		fmt.Sprintf("option go_package=\"go_agent/telemetry/genproto/ros/%q\";", g.localRosPackage),
	}
}

func (g *genProtoIDL) FileName() string { return g.OutputFilename + ".proto" }
func (g *genProtoIDL) FileType() string { return "protoidl" }

func (g *genProtoIDL) Namers(c *Context) NameSystems {
	return NameSystems{
		"local": localNamer{g.localPackage},
	}
}

func (g *genProtoIDL) Imports(c *Context) (imports []string) {
	lines := []string{}

	// for _, line := range g.imports.ImportLines() {
	// 	lines = append(lines, line)
	// }
	//
	return lines
}

func (g *genProtoIDL) GenerateType(c *Context, t *Type, w io.Writer) error {
	sw := NewSnippetWriter(w, c, "$", "$")
	b := bodyGen{
		locator: &protobufLocator{
			namer:           c.Namers["proto"].(ProtobufFromGoNamer),
			universe:        c.Universe,
			localRosPackage: g.localRosPackage.String(),
		},
		localPackage: g.localPackage,
		t:            t,
	}
	switch t.Kind {
	case RosMsgType:
		return b.doRosMsgType(sw)
	default:
		return b.unknown(sw)
	}
}

type bodyGen struct {
	locator      ProtobufLocator
	localPackage Name
	t            *Type
}

func (b bodyGen) doRosMsgType(sw *SnippetWriter) error {
	if len(b.t.Name.Name) == 0 {
		return nil
	}

	var alias *Type
	var fields []protoField

	if alias == nil {
		alias = b.t
	}
	imports, err := getImports(b.locator, alias, b.localPackage)
	slices.Sort(imports)
	imports = slices.Compact(imports)
	if err != nil {
		return fmt.Errorf("Failed to get imports for %v: %v", b.t, err)
	}

	if fields == nil {
		memberFields, err := memberToFields(b.locator, alias, b.localPackage)
		if err != nil {
			return fmt.Errorf("type %v cannot be converted to protobuf: %v", b.t, err)
		}
		fields = memberFields
	}

	out := sw.Out()
	for _, im := range imports {
		p := filepath.Base(im)

		if strings.ToLower(p) == strings.ToLower(b.t.Name.Name) {
			continue // avoid self importing
		}
		fmt.Fprintf(out, "import \"%s.proto\";\n", im)
	}
	fmt.Fprintf(out, "\n\n")
	sw.Do(`message $.Name.Name$`, b.t)
	fmt.Fprintf(out, "\n{\n")
	for i, field := range fields {
		fmt.Fprintf(out, " ")

		if field.Repeated {
			fmt.Fprintf(out, "repeated ")
		}

		sw.Do(`$.Type|local$ $.Name$ = $.Tag$`, field)
		fmt.Fprintf(out, ";\n")
		if i != len(fields)-1 {
			fmt.Fprintf(out, "\n")
		}
	}
	fmt.Fprintf(out, "}\n\n")
	return nil
}
func (b bodyGen) unknown(sw *SnippetWriter) error {
	return fmt.Errorf("cannot generate %#v", b.t)
}

type ProtobufFromGoNamer interface {
	RosNameToProtoName(name Name) Name
}

type ProtobufLocator interface {
	ProtoTypeFor(t *Type) (*Type, error)
	RosTypeForName(name Name) *Type
	CastTypeName(name Name) string
	GetMessageDef(name Name, msgName string) (*MessageDefinition, error)
	RosPackageName() string
}

type protobufLocator struct {
	namer    ProtobufFromGoNamer
	tracker  ImportTracker
	universe Universe

	localRosPackage string
}

func (p protobufLocator) CastTypeName(name Name) string {
	if name.Package == p.localRosPackage {
		return name.Name
	}

	return name.String()
}

func (p protobufLocator) RosPackageName() string {
	return p.localRosPackage
}

func (p protobufLocator) GetMessageDef(name Name, msgName string) (*MessageDefinition, error) {
	if t, ok := p.universe[name.Path]; ok {
		if msg, ok := t.MessageDefs[name.Path]; ok {
			ret := msg[msgName]
			// log.Printf("Message definiton: %s\n", msgName)
			// for i, def := range ret.Definitions {
			// 	log.Printf("Def %d: %v\n", i, def)
			// }
			return ret, nil
		}
	}
	log.Printf("Message type %v not found in package %v", msgName, name.Path)
	return nil, fmt.Errorf("Message type %v not found in package %v", msgName, name.Package)
}

func (p protobufLocator) RosTypeForName(name Name) *Type {
	if len(name.Package) == 0 {
		name.Package = p.localRosPackage
	}
	return p.universe.Type(name)
}

func (p protobufLocator) ProtoTypeFor(t *Type) (*Type, error) {
	if t.Kind == Protobuf {
		return t, nil
	}

	if t.Kind == Builtin {
		return t, nil
	}

	if t, ok := isFundamentalPrototype(t); ok {
		return t, nil
	}

	if t.Kind == RosMsgType {
		t := &Type{
			Name: p.namer.RosNameToProtoName(t.Name),
			Kind: Protobuf,
		}
		return t, nil
	}

	return nil, errUnrecognizedType
}

type protoField struct {
	LocalPackage Name
	Tag          int
	Name         string
	Type         *Type
	Repeated     bool
}

func isFundamentalPrototype(t *Type) (*Type, bool) {
	switch t.Kind {
	case Builtin:
		switch t.Name.Name {
		case "string", "uint32", "int32", "uint64", "int64", "bool":
			return &Type{Name: Name{Name: t.Name.Name}, Kind: Protobuf}, true
		case "int":
			return &Type{Name: Name{Name: "int64"}, Kind: Protobuf}, true
		case "uint":
			return &Type{Name: Name{Name: "uint64"}, Kind: Protobuf}, true
		case "float64":
			return &Type{Name: Name{Name: "double"}, Kind: Protobuf}, true
		case "float32", "float":
			return &Type{Name: Name{Name: "float"}, Kind: Protobuf}, true
		}
	}
	return t, false
}

func memberTypeToProtobufField(locator ProtobufLocator, field *protoField, t *Type) error {
	var err error
	switch t.Kind {
	case Builtin:
		field.Type, err = locator.ProtoTypeFor(t)
	case Protobuf:
		field.Type, err = locator.ProtoTypeFor(t)
	case RosMsgType:
		if len(t.Name.Name) == 0 {
			return errUnrecognizedType
		}
		field.Type, err = locator.ProtoTypeFor(t)
	default:
		return errUnrecognizedType
	}
	return err
}
func getImports(locator ProtobufLocator, t *Type, localPackage Name) ([]string, error) {
	imports := []string{}
	msgDef, err := locator.GetMessageDef(localPackage, t.String())

	if err != nil {
		return imports, err
	}

	for _, f := range msgDef.Fields {
		if f.TypePkg.Name != "" {
			imports = append(imports, fmt.Sprintf("%s/%s", f.TypePkg.Name, f.Type.Name))
		} else if !f.Type.isBuiltin() {
			log.Printf("Not a builtin type %s and its package is %s", f.Type.Name.Name, msgDef.RosPkgName.Name)
			imports = append(imports, fmt.Sprintf("%s/%s", msgDef.RosPkgName.Name, f.Type.Name.Name))
		}
	}
	return imports, nil
}

func memberToFields(locator ProtobufLocator, t *Type, localPackage Name) ([]protoField, error) {
	fields := []protoField{}
	msgDef, err := locator.GetMessageDef(localPackage, t.String())
	if err != nil {
		return fields, err
	}

	for i, f := range msgDef.Fields {
		field := protoField{
			LocalPackage: localPackage,
			Tag:          -1,
		}

		//   if f.TypePkg.Name != "" {
		// 	continue // this is an import statement
		// }
		for k := range msgDef.Imports {
			if strings.Contains(k, f.Type.Name.Name) {
				m := strings.Split(k, "/")
				f.Type.Name.Name = fmt.Sprintf("%s.%s", m[0], f.Type.Name.Name)
			}
		}

		if field.Type == nil {
			if err := memberTypeToProtobufField(locator, &field, &f.Type); err != nil {
				f.Type.Kind = RosMsgType
				//prepend typename with proto package name
				// field.Name = f.Type.Name.Name
				nerr := memberTypeToProtobufField(locator, &field, &f.Type)
				if nerr != nil {
					log.Printf("unable to embed type %q as field %q in %q: %v", &f.Type, f.Name, field.Type, err)
					return nil, fmt.Errorf("unable to embed type %q as field %q in %q: %v", &f.Type, f.Name, field.Type, err)
				} else {
					log.Printf("Embedding type %q as field %q in %q", &f.Type, f.Name, field.Name)
				}
			}
		}

		if len(field.Name) == 0 {
			field.Name = f.Name
		}

		if f.TypeArray {
			field.Repeated = true
		}

		field.Tag = i + 1
		fields = append(fields, field)
	}

	return fields, nil
}

func assembleProtoFile(w io.Writer, f *File) {
	w.Write(f.Header)

	fmt.Fprintf(w, "syntax = \"proto3\";\n\n")

	if len(f.PackageName) != 0 {
		fmt.Fprintf(w, "package %s;\n", f.PackagePath)
		fmt.Fprintf(w, "option go_package = \"go_agent/telemetry/genproto/ros/%s\";\n\n", f.PackagePath)
	}
	w.Write(f.Body.Bytes())
}

func NewProtoFile() *DefaultFileType {
	return &DefaultFileType{
		Assemble: assembleProtoFile,
	}
}
