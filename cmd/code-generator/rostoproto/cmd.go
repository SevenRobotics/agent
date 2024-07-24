package rostoproto

import (
	"fmt"
	"go_agent/cmd/code-generator/rostoproto/util"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	flag "github.com/spf13/pflag"
)

type GeneratorUtil struct {
	RosMsgFiles   string
	RosPackages   string
	OutputDir     string
	ProtoImport   []string
	Conditional   string
	Clean         bool
	OnlyIDL       bool
	KeepGogoproto bool
}

func NewGen() *GeneratorUtil {
	defaultOutputDir := "../../../telemetry/protobuf/"
	return &GeneratorUtil{
		OutputDir:   defaultOutputDir,
		RosPackages: "",
	}
}

func (g *GeneratorUtil) BindFlags(flag *flag.FlagSet) {
	flag.StringVarP(&g.RosMsgFiles, "ros-msg-file", "m", "", "Ros msg file to parse and generate proto for")
	flag.StringVarP(&g.RosPackages, "ros-package", "p", g.RosPackages, "comma separated list of ros packages to get msg files from. Directories prefixed with '-' are not generated, directories prefixed with '+' only create types with explicit IDL instructions")
	flag.StringVar(&g.OutputDir, "output-dir", g.OutputDir, "the default output dir for generated proto files")
	flag.StringSliceVar(&g.ProtoImport, "proto-import", g.ProtoImport, "A search path for imported protobufs")
}

func Run(g *GeneratorUtil) {
	p := NewWithOptions(Options{BuildTags: []string{"proto"}})
	var allPackages []string
	var allMsgs []string
	if len(g.RosMsgFiles) != 0 {
		allMsgs = append(allMsgs, strings.Split(g.RosMsgFiles, ",")...)
	}

	if len(g.RosPackages) != 0 {
		allPackages = append(allPackages, strings.Split(g.RosPackages, ",")...)
	}
	rospkgs, err := util.FindRosPackages()
	if err != nil {
		log.Fatal("Could not find ros packages %v", err)
	}

	if len(rospkgs) > 0 {
		allPackages = append(allPackages, rospkgs...)
	}

	//sort and remove any duplicates
	slices.Sort(allPackages)
	slices.Compact(allPackages)

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

		// remove packages which do not have msg types
		if _, err := os.Stat(d + "/msg"); os.IsNotExist(err) {
			modifier.output = false
			inputPackageModifiers[d] = modifier
			continue
		}

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

	loadablePackages := make([]string, 0, len(allPackages))

	for _, p := range allPackages {
		if o, ok := inputPackageModifiers[p]; ok {
			if !o.output {
				continue
			}
		}
		loadablePackages = append(loadablePackages, p)
	}

	if err := p.LoadPackages(loadablePackages...); err != nil {
		log.Fatalf("Failed to load packages %v", err)
	}

	c, err := NewContext(
		p,
		NameSystems{
			"public": NewNamer(3),
		},
	)

	if err != nil {
		log.Fatalf("Failed creating a new context %v", err)
	}

	c.FileTypes["protoidl"] = NewProtoFile()

	protoBufNames := NewProtobufNamer()
	outputPackages := []Target{}

	for _, input := range c.Inputs {
		pkg := c.Universe[filepath.Base(input)]
		if pkg == nil {
			fmt.Printf("pkg for input %v is empty", input)
		}
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		dir := filepath.Join(basepath, "../../../telemetry/protobuf/", pkg.Name)
		for _, msg := range pkg.MessageDefs[pkg.Name] {
			protopkg := newProtobufPackage(pkg.Path, dir, msg.Name.Name, true)
			protoBufNames.Add(protopkg)
			outputPackages = append(outputPackages, protopkg)
		}
	}

	c.Namers["proto"] = protoBufNames

	if err := c.ExecuteTargets(outputPackages); err != nil {
		log.Fatalf("Failed executing local generator: %v", err)
	}

	return
}
