package rostoproto

import (
	"bytes"
	"fmt"
	"text/template"
)

var tpl = template.Must(template.New("").Parse(
	`func Convert{{ .Name.Name }}(rosMsg *gengo_{{ .RosPkgName.Name }}.{{ .Name.Name }}) (proto_{{ .RosPkgName.Name }}.{{ .Name.Name }}, error) {
    ret := proto_{{ .RosPkgName.Name }}.{{ .Name.Name }}{}
    if rosMsg == nil {
      return ret, fmt.Errorf("Cannot not convert nil msg")
    }
    
    {{range .Fields}}
      {{ if eq .Type.Kind "Builtin"}}
        ret.{{ .Name }} = rosMsg.{{ .Name }}
      {{ else }}
        ret.{{ .Name }}, _ = Convert{{ .Type.Name.Name}}(rosMsg.{{ .Name }})   
      {{ end }}
    {{end}}
    
    return ret, nil 
  }
`))

var tpl_imports = template.Must(template.New("").Parse(
	`gengo_{{ .RosPkgName.Name }} "go_agent/telemetry/gengo/ros/{{ .RosPkgName.Name }}"
  proto_{{ .RosPkgName.Name}} "go_agent/telemetry/genproto/ros/{{ .RosPkgName.Name }}"
  `))

func WriteHeader() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "package converter\n\n")
	fmt.Fprintf(&buf, "import (\n")
	fmt.Fprintf(&buf, "\"fmt\"\n")
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
