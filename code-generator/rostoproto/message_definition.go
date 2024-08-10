package rostoproto

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func firstCharToUpper(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func snakeToCamel(in string) string {
	tmp := []rune(in)
	tmp[0] = unicode.ToUpper(tmp[0])

	for i := 0; i < (len(tmp) - 1); i++ {
		if tmp[i] == '_' {
			tmp[i+1] = unicode.ToUpper(tmp[i+1])
			tmp = append(tmp[:i], tmp[i+1:]...)
			i--
		}
	}

	return string(tmp)
}

func camelToSnake(in string) string {
	tmp := []rune(in)
	tmp[0] = unicode.ToLower(tmp[0])

	for i := 0; i < len(tmp); i++ {
		if unicode.IsUpper(tmp[i]) {
			tmp[i] = unicode.ToLower(tmp[i])
			tmp = append(tmp[:i], append([]rune{'_'}, tmp[i:]...)...)
		}
	}

	return string(tmp)
}

// Ros message type definition
type Definition struct {
	RosType   Type
	ProtoType Type
	Name      string
	Value     string
	HasQuotes bool
}

// Ros Message field
type Field struct {
	TypeArray    bool // is this an array type
	TypePkg      Name
	Type         Type
	Name         string
	NameOverride string
}

// Ros Message Definition
type MessageDefinition struct {
	RosPkgName     Name
	Name           Name
	Fields         []Field
	Definitions    []Definition
	DefinitionsStr string
	Imports        map[string]struct{}
	Type           Type
}

func parseField(rosPkgName string, res *MessageDefinition, typ string, name string) {
	f := Field{}

	f.Name = snakeToCamel(name)
	// if backward conversion to snake not possible, save the name for override
	if camelToSnake(f.Name) != name {
		f.NameOverride = name
	}

	//split an array and its type
	m := regexp.MustCompile(`^(.+?)(\[.*?\])$`).FindStringSubmatch(typ)
	if m != nil {
		f.TypeArray = true
		f.Type = Type{Name: Name{Name: m[1]}}
		if _, ok := Builtins.Types[f.Type.Name.Name]; !ok {
			f.Type.Kind = RosMsgType
		} else {
			f.Type.Kind = Builtin
		}
	} else {
		f.Type = Type{Name: Name{Name: typ}}
		if _, ok := Builtins.Types[f.Type.Name.Name]; !ok {
			f.Type.Kind = RosMsgType
		} else {
			f.Type.Kind = Builtin
		}
	}

	//check if this is a package import field
	f.TypePkg, f.Type = func() (Name, Type) {
		parts := strings.Split(f.Type.Name.Name, "/")
		if len(parts) == 2 {
			return Name{Name: parts[0]}, Type{Name: Name{Name: parts[1]}}
		}

		switch f.Type.Name.Name {
		case "Bool", "ColorRGBA", "Empty", "Float32MultiArray", "Float32",
			"Float64MultiArray", "Float64", "Header", "Int8MultiArray",
			"Int8", "Int16MultiArray", "Int16", "Int32MultiArray", "Int32",
			"Int64MultiArray", "Int64", "MultiArrayDimension", "MultiarrayLayout",
			"String", "UInt8MultiArray", "UInt8", "UInt16MultiArray", "UInt16",
			"UInt32MultiArray", "UInt32", "UInt64MultiArray", "UInt64":
			return Name{Name: "std_msgs"}, Type{Name: Name{Name: parts[0]}, Kind: RosMsgType}

		case "bool", "string":
			return Name{}, f.Type

		case "int8":
			return Name{}, Type{Name: Name{Name: "int32"}, Kind: Builtin, Original: "int8"}

		case "int16":
			return Name{}, Type{Name: Name{Name: "int32"}, Kind: Builtin, Original: "int16"}

		case "int32":
			return Name{}, *Builtins.Types["int32"]

		case "uint8":
			return Name{}, Type{Name: Name{Name: "uint32"}, Kind: Builtin, Original: "uint8"}

		case "uint16":
			return Name{}, Type{Name: Name{Name: "uint32"}, Kind: Builtin, Original: "uint16"}

		case "uint32":
			return Name{}, *Builtins.Types["uint32"]

		case "int64":
			return Name{}, f.Type

		case "uint64":
			return Name{}, f.Type

		case "float32":
			return Name{}, *Builtins.Types["float"]

		case "float64":
			return Name{}, *Builtins.Types["double"]

		case "time", "duration":
			return Name{}, Type{Name: Name{Name: firstCharToUpper(f.Type.Name.Name)}, Kind: Protobuf}

		case "byte", "char":
			return Name{}, Type{Name: Name{Name: "uint32"}, Kind: Builtin, Original: "byte"}
		}
		return Name{}, f.Type
	}()

	res.Fields = append(res.Fields, f)
}

func parseDefinition(res *MessageDefinition, typ string, name string, val string) {
	d := Definition{
		RosType: Type{Name: Name{Name: typ}},
		Name:    name,
		Value:   val,
	}

	d.Value = strings.ReplaceAll(d.Value, "\"", "\\\"")

	d.ProtoType = func() Type {
		switch d.RosType.Name.Name {
		case "byte":
			return *Builtins.Types["int8"]
		case "char":
			return *Builtins.Types["uint8"]
		case "float32":
			return *Builtins.Types["float"]
		case "float64":
			return *Builtins.Types["double"]

		}
		return d.RosType
	}()
	res.Definitions = append(res.Definitions, d)
}

func ParseMessageDefinition(rosPkgName string, name string, content string) (*MessageDefinition, error) {
	res := &MessageDefinition{
		RosPkgName: Name{Name: rosPkgName},
		Name:       Name{Name: firstCharToUpper(name)},
	}

	for _, line := range strings.Split(content, "\n") {
		//remove all comments
		line = regexp.MustCompile("#.*$").ReplaceAllString(line, "")

		//remove leading and trailing whitespace
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		i := strings.IndexAny(line, " \t")
		if i < 0 {
			return nil, fmt.Errorf("unable to parse line (%s)", line)
		}

		var typ string
		typ, line = line[:i], line[i+1:]

		line = strings.TrimLeft(line, " \t")

		i = strings.IndexByte(line, '=')
		if i < 0 {
			name := line
			parseField(rosPkgName, res, typ, name)
		} else {
			name, val := line[:i], line[i+1:]
			name = strings.TrimLeft(name, " \t")
			val = strings.TrimLeft(val, " \t")
			parseDefinition(res, typ, name, val)
		}
	}

	res.DefinitionsStr = func() string {
		var tmp []string
		for _, d := range res.Definitions {
			tmp = append(tmp, d.RosType.Name.Name+" "+d.Name+"="+d.Value)
		}
		return strings.Join(tmp, ",")
	}()

	res.Imports = map[string]struct{}{}

	for _, f := range res.Fields {
		switch f.TypePkg.Name {
		case "":
			//do nothing not an import field
		default:
			res.Imports[f.TypePkg.Name+"/"+f.Type.Name.Name] = struct{}{}
		}
	}
	res.Type = Type{Name: Name{Name: name}, Kind: RosMsgType}
	return res, nil
}
