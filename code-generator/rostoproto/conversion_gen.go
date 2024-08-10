package rostoproto

import (
	"bytes"
	"fmt"
	"text/template"
)

var tpl = template.Must(template.New("").Parse(
	`func Convert{{ .Name.Name }}(rosMsg gengo_{{ .RosPkgName.Name }}.{{ .Name.Name }}) (proto_{{ .RosPkgName.Name }}.{{ .Name.Name }}, error) {
    
    ret := proto_{{ .RosPkgName.Name }}.{{ .Name.Name }}{}
    
    //if rosMsg == nil {
    //  return ret, fmt.Errorf("Cannot not convert {{ .Name.Name }}, msg is nil")
    //}
    
    {{ range .Fields }}
      {{if and (ne .Type.Kind "Builtin") (and (ne .Type.Name.Name "Duration") (ne .Type.Name.Name "Time"))}}
      var err error
      {{ break }} 
      {{end}}
    {{end}}

    {{range $i, $a := .Fields}}
      {{ if eq .Type.Kind "Builtin"}}
        {{ if .TypeArray}}
          {{ if or (eq .Type.Original "int8") (eq .Type.Original "int16") }}
          t{{$i}} := make([]int32,0,len(rosMsg.{{ .Name }}))
          for _, m := range rosMsg.{{ .Name }} {
            a := int32(m)
            t{{$i}} = append(t{{$i}}, a)
          }
          {{ else if or (eq .Type.Original "uint8") (eq .Type.Original "uint16") (eq .Type.Original "byte")}}
          t{{$i}} := make([]uint32,0,len(rosMsg.{{ .Name }}))
          for _, m := range rosMsg.{{ .Name }} {
            a := uint32(m)
            t{{$i}} = append(t{{$i}}, a)
          }
          {{ else if eq .Type.Name.Name "double"}}
          t{{$i}} := make([]float64,0,len(rosMsg.{{ .Name }}))
          for _, m := range rosMsg.{{ .Name }} {
            t{{$i}} = append(t{{$i}}, m)
          }
          {{ else if eq .Type.Name.Name "float"}}
          t{{$i}} := make([]float32,0,len(rosMsg.{{ .Name }}))
          for _, m := range rosMsg.{{ .Name }} {
            t{{$i}} = append(t{{$i}}, m)
          }
          {{ else }}
          t{{$i}} := make([]{{ .Type.Name.Name }},0,len(rosMsg.{{ .Name }}))
          for _, m := range rosMsg.{{ .Name }} {
            t{{$i}} = append(t{{$i}}, m)
          }
          {{end}}
          ret.{{ .Name }} = append(ret.{{ .Name }}, t{{$i}}...)
        {{else}}
          {{ if or (eq .Type.Original "uint8") (eq .Type.Original "uint16") (eq .Type.Original "byte")}}
          ret.{{ .Name }} = uint32(rosMsg.{{ .Name }})
          
          {{ else if or (eq .Type.Original "int8") (eq .Type.Original "int16") }}
          ret.{{ .Name }} = int32(rosMsg.{{ .Name }})

          {{ else }}
          ret.{{ .Name }} = rosMsg.{{ .Name }}
          
          {{end}}
        {{end}}
      {{ else }}
        {{if or (eq .Type.Name.Name "Duration") (eq .Type.Name.Name "Time") }}
          {{ continue }}
        {{ else }}
        {{ if .TypeArray}}
          for _, m := range rosMsg.{{ .Name }} {
            c , err := Convert{{.Type.Name.Name}}(m)
            if err != nil {
              return ret, err 
            }
            ret.{{ .Name }} = append(ret.{{ .Name }}, &c)
          }
        {{else}}
        t{{$i}} , err  := Convert{{ .Type.Name.Name}}(rosMsg.{{ .Name }})
        ret.{{ .Name }} = &t{{$i}}
        {{end}}
        {{ end }}
      if err != nil {
        return ret, err 
      }

      {{ end }}
    {{end}}
    
    return ret, nil 
  }
`))

var tpl_imports = template.Must(template.New("").Parse(
	`gengo_{{ .RosPkgName.Name }} "go_agent/telemetry/gengo/ros/{{ .RosPkgName.Name }}"
  proto_{{ .RosPkgName.Name}} "go_agent/telemetry/genproto/ros/{{ .RosPkgName.Name }}"
  `))

var tpl_switch = template.Must(template.New("").Parse(
	`case "{{.Name.Name}}":
return &channel.BuilderUtil[gengo_{{ .RosPkgName.Name }}.{{ .Name.Name }},proto_{{.RosPkgName.Name}}.{{.Name}}]{
  MsgConverter: Convert{{.Name.Name}},
  Serializer: Serialize{{.Name.Name}},
}, nil
  `))

var tpl_serialize = template.Must(template.New("").Parse(
	`func Serialize{{ .Name.Name }}(msg proto_{{ .RosPkgName.Name }}.{{ .Name.Name }}) ([]byte, error) {
    return proto.Marshal(&msg)
  }`))

var tpl_init = template.Must(template.New("").Parse(
	`func AssignBuilder() bool {
	utils.NewBuilder("ros-rmq", GetBuilderFromName)
	return true
  }`))

func WriteHeader() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "package converter\n\n")
	fmt.Fprintf(&buf, "import (\n")
	fmt.Fprintf(&buf, "\"fmt\"\n")
	fmt.Fprintf(&buf, "\"github.com/golang/protobuf/proto\"\n")
	fmt.Fprintf(&buf, "\"go_agent/iface\"\n")
	fmt.Fprintf(&buf, "\"go_agent/utils\"\n")
	fmt.Fprintf(&buf, "\"go_agent/telemetry/cmd/channel\"\n")
	return buf.String()
}

func CloseImports() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, ")\n")
	return buf.String()
}

func WriteImports(res *MessageDefinition) (string, error) {
	var buf bytes.Buffer
	err := tpl_imports.Execute(&buf, res)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func WriteConverter(res *MessageDefinition) (string, error) {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, res)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func WriteSerializer(res *MessageDefinition) (string, error) {
	var buf bytes.Buffer
	err := tpl_serialize.Execute(&buf, res)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func WriteSwitch(res *MessageDefinition) (string, error) {
	var buf bytes.Buffer
	err := tpl_switch.Execute(&buf, res)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func WriteGetFuncBegin() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "func GetBuilderFromName(name string) (iface.Builder, error) {\n")
	fmt.Fprintf(&buf, "switch name {\n")
	return buf.String()
}

func WriteGetFuncEnd() string {
	var buf, initBuf bytes.Buffer
	fmt.Fprintf(&buf, "default:\n")
	fmt.Fprintf(&buf, "return nil,fmt.Errorf(\"Unrecognized Type\")\n")
	fmt.Fprintf(&buf, "}\n}\n")
	err := tpl_init.Execute(&initBuf, nil)
	if err != nil {
		return ""
	}
	return buf.String() + initBuf.String()
}
