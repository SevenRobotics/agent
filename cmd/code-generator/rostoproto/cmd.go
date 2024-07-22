package rostoproto

import (
	"log"
	"strings"

	flag "github.com/spf13/pflag"

	"k8s.io/gengo/v2/parser"
)

type Generator struct {
	RosMsgFiles   string
	RosPackages   string
	OutputDir     string
	ProtoImport   []string
	Conditional   string
	Clean         bool
	OnlyIDL       bool
	KeepGogoproto bool
}

func NewGen() *Generator {
	defaultOutputDir := "../../../telemetry/protobuf/"
	return &Generator{
		OutputDir:   defaultOutputDir,
		RosPackages: "",
	}
}

func (g *Generator) BindFlags(flag *flag.FlagSet) {
	flag.StringVarP(&g.RosMsgFiles, "ros-msg-file", "m", "", "Ros msg file to parse and generate proto for")
	flag.StringVarP(&g.RosPackages, "ros-package", "p", g.RosPackages, "comma separated list of ros packages to get msg files from. Directories prefixed with '-' are not generated, directories prefixed with '+' only create types with explicit IDL instructions")
	flag.StringVar(&g.OutputDir, "output-dir", g.OutputDir, "the default output dir for generated proto files")
	flag.StringSliceVar(&g.ProtoImport, "proto-import", g.ProtoImport, "A search path for imported protobufs")
	flag.StringVar(&g.Conditional, "conditional", g.Conditional, "An optional golang build tag condition to add to the generated Go Code")
	flag.BoolVar(&g.Clean, "clean", g.Clean, "If true, remove all generated files for the specified Packages.")
	flag.BoolVar(&g.OnlyIDL, "only-idl", g.OnlyIDL, "If true, only generate the IDL for each package.")
	flag.BoolVar(&g.KeepGogoproto, "keep-gogoproto", g.KeepGogoproto, "If true, the generated IDL will contain gogoprotobuf extensions which are normally removed")
}

func Run(g *Generator) {
	p := parser.NewWithOptions(parser.Options{BuildTags: []string{"proto"}})
	var allPackages []string
	var allMsgs []string
	if len(g.RosMsgFiles) != 0 {
		allMsgs = append(allMsgs, strings.Split(g.RosMsgFiles, ",")...)
	}

	if len(g.RosPackages) != 0 {
		allPackages = append(allPackages, strings.Split(g.RosPackages, ",")...)
	}

	if len(allPackages) == 0 && len(allMsgs) == 0 {
		log.Fatalf("No Ros packages or Ros msg files received. Exiting")
	}

	type modifier struct {
		allTypes bool
		output   bool
		name     string
	}

	inputPackageModifiers := map[string]modifier{}
	inputMsgModifiers := map[string]modifier{}
	packages := make([]string, 0, len(allPackages))
	messages := make([]string, 0, len(allMsgs))

	for _, d := range allPackages {
		modifier := modifier{allTypes: true, output: true}

		switch {
		case strings.HasPrefix(d, "-"):
			d = d[1:]
			modifier.output = false
		case strings.HasPrefix(d, "+"):
			d = d[1:]
			modifier.allTypes = false

			name := protoSafePackage(d)
			modifier.name = name
			packages = append(packages, d)
			inputPackageModifiers[d] = modifier
		}
	}

	for _, d := range allMsgs {
		modifier := modifier{allTypes: true, output: true}
		switch {
		case strings.HasPrefix(d, "-"):
			d = d[1:]
			modifier.output = false
		case strings.HasPrefix(d, "+"):
			d = d[1:]
			modifier.allTypes = false

			name := protoSafePackageMsg(d)
			modifier.name = name
			messages = append(messages, d)
			inputMsgModifiers[d] = modifier
		}
	}

}
