package rostoproto

import (
	"bytes"
	"fmt"
	"go_agent/cmd/code-generator/rostogo"
	"go_agent/cmd/code-generator/rostoproto/util"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"

	flag "github.com/spf13/pflag"
)

type GeneratorUtil struct {
	RosMsgFiles     string
	RosPackages     string
	OutputDir       string
	ProtoDir        string
	GoDir           string
	RosSubOutputDir string
	Blacklist       []string
	ProtoImport     []string
	Conditional     string
	Clean           bool
	OnlyIDL         bool
	KeepGogoproto   bool
}

type GeneratorState struct {
	//map of package name to its Messages: pkg -> [msgName -> MessageDefinition]
	RosMsgPkgs map[string]map[string]*MessageDefinition
	// map of ros messages and their proto paths
	ProtoPkgPath map[string]string
	// map of ros messages and their generared proto paths
	ProtoImportPath map[string]string
}

func NewGen() *GeneratorUtil {
	defaultOutputDir := "../../../telemetry/genproto/ros/"
	relativeProtoDir := "../../../telemetry/protobuf/ros/"
	rosSubOutputDir := "../../../ros_subscribers/ros/"
	goOutputDir := "../../../telemetry/gengo/ros/"
	blacklist := []string{"turtlesim", "tf"}
	return &GeneratorUtil{
		OutputDir:       defaultOutputDir,
		ProtoDir:        relativeProtoDir,
		GoDir:           goOutputDir,
		RosSubOutputDir: rosSubOutputDir,
		Blacklist:       blacklist,
		RosPackages:     "",
	}
}

func (g *GeneratorUtil) BindFlags(flag *flag.FlagSet) {
	flag.StringVarP(&g.RosMsgFiles, "ros-msg-file", "m", "", "Ros msg file to parse and generate proto for")
	flag.StringVarP(&g.RosPackages, "ros-package", "p", g.RosPackages, "comma separated list of ros packages to get msg files from. Directories prefixed with '-' are not generated, directories prefixed with '+' only create types with explicit IDL instructions")
	flag.StringVar(&g.OutputDir, "output-dir", g.OutputDir, "the default output dir for generated proto files")
	flag.StringSliceVar(&g.ProtoImport, "proto-import", g.ProtoImport, "A search path for imported protobufs")
}

func Run(g *GeneratorUtil) {
	genState := GeneratorState{}
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
		log.Fatalf("Could not find ros packages %v", err)
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

		for _, b := range g.Blacklist {
			if filepath.Base(d) == b {
				modifier.output = false
				inputPackageModifiers[d] = modifier
				continue
			}
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
	genState.RosMsgPkgs = map[string]map[string]*MessageDefinition{}
	genState.ProtoPkgPath = map[string]string{}
	genState.ProtoImportPath = map[string]string{}

	for _, p := range allPackages {
		if o, ok := inputPackageModifiers[p]; ok {
			if !o.output {
				continue
			}
		}
		loadablePackages = append(loadablePackages, p)
	}
	// first lets create go message types for use with goroslib subscribers

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
	c.FileTypes["go"] = NewRosSubFile()

	protoBufNames := NewProtobufNamer()
	RosSubNames := NewRosSubNamer()
	outputPackages := []Target{}
	RosSubPackages := []Target{}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	converterDir := filepath.Join(basepath, filepath.Dir(g.GoDir), "converter")
	err = os.MkdirAll(converterDir, 0755)
	if err != nil {
		log.Fatalf("failed to create dir %s: %v", converterDir, err)
	}
	converterGoFile := converterDir + "/converter.go"
	f, err := os.Create(converterGoFile)
	if err != nil {
		log.Fatalf("Failed to create file %s : %v", converterGoFile, err)
	}

	f.Write([]byte(WriteHeader()))
	var importBuf bytes.Buffer
	var converterBuf bytes.Buffer

	for _, input := range c.Inputs {
		pkg := c.Universe[filepath.Base(input)]
		if _, ok := genState.RosMsgPkgs[pkg.Name]; !ok {
			genState.RosMsgPkgs[pkg.Name] = map[string]*MessageDefinition{}
		}
		if pkg == nil {
			fmt.Printf("pkg for input %v is empty", input)
		}

		dir := filepath.Join(basepath, g.ProtoDir, pkg.Name)
		goDir := filepath.Join(basepath, g.GoDir, pkg.Name)
		subDir := filepath.Join(basepath, g.RosSubOutputDir, pkg.Name)

		err = rostogo.ImportPackage(input, goDir)
		if err != nil {
			log.Fatalf("Failed to import ros package %s: %v", pkg.Name, err)
		}

		//remove all generated action files
		filepath.WalkDir(goDir, func(path string, d fs.DirEntry, err error) error {
			if strings.HasPrefix(d.Name(), "action") && filepath.Ext(d.Name()) == ".go" {
				err := os.Remove(path)
				if err != nil {
					log.Fatalf("Failed to remove file %s: %v", path, err)
				}
			}
			return nil
		})

		// generate package.go files for ros_subscribers
		err := GenPackageFile(pkg.Name+"_sub", subDir)

		if err != nil {
			log.Fatalf("failed to create package.go for %s : %v", pkg.Name, err)
		}

		t := genState.RosMsgPkgs[pkg.Name]

		for _, msg := range pkg.MessageDefs[pkg.Name] {
			imp, _ := WriteImports(msg)
			importBuf.Write([]byte(imp))
			break
		}

		for _, msg := range pkg.MessageDefs[pkg.Name] {
			con, _ := WriteConverter(msg)
			con = con + "\n\n"
			converterBuf.Write([]byte(con))
			protopkg := newProtobufPackage(pkg.Path, dir, msg.Name.Name, true)
			protoBufNames.Add(protopkg)
			outputPackages = append(outputPackages, protopkg)
			subpkg := NewRosSubPackage(pkg.Path, subDir, msg.Name.Name, true)
			RosSubNames.Add(subpkg)
			RosSubPackages = append(RosSubPackages, subpkg)
			t[msg.Name.Name] = msg
			genState.ProtoPkgPath[msg.Name.Name] = filepath.Join(dir, msg.Name.Name) + ".proto"
		}
		genState.RosMsgPkgs[pkg.Name] = t
	}

	importBuf.Write([]byte(CloseImports()))
	f.Write(importBuf.Bytes())
	f.Write(converterBuf.Bytes())

	cmd := exec.Command("gofmt", "-s", "-w", converterGoFile)
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		log.Print(string(out))
	}
	if err != nil {
		log.Println(strings.Join(cmd.Args, " "))
		log.Fatalf("unable to apply gofmt to %s: %v", converterGoFile, err)
	}

	c.Namers["proto"] = protoBufNames
	c.Namers["go"] = RosSubNames
	deps := deps(c, protoBufNames.packages)

	order, err := importOrder(deps)

	if err != nil {
		log.Fatalf("Failed to order packages by imports: %v", err)
	}

	topologicalPos := map[string]int{}

	for i, p := range order {
		topologicalPos[p] = i
	}

	sort.Sort(positionOrder{topologicalPos, protoBufNames.packages})
	var sortedOutputPackages []Target

	for _, protoPkg := range protoBufNames.packages {
		sortedOutputPackages = append(sortedOutputPackages, protoPkg)
	}

	if err := c.ExecuteTargets(outputPackages); err != nil {
		log.Fatalf("Failed executing local generator: %v", err)
	}

	if err := c.ExecuteTargets(RosSubPackages); err != nil {
		log.Fatalf("Failed executing ros sub generator: %v", err)
	}

	if _, err := exec.LookPath("protoc"); err != nil {
		log.Fatalf("unable to find protoc: %v", err)
	}

	_, b, _, _ = runtime.Caller(0)
	basepath = filepath.Dir(b)
	search_args := []string{}
	for _, outputPackage := range outputPackages {
		p := outputPackage.(*ProtobufPackage)

		path := filepath.Join(basepath, g.ProtoDir, p.ImportPath())
		tmp := strings.Split(path, "/")
		includePath := strings.Join(tmp[:len(tmp)-2], "/")
		search_args = append(search_args, fmt.Sprintf("--proto_path=%s", includePath))
		outputPath := filepath.Join(basepath, g.OutputDir)
		tmp = strings.Split(outputPath, "/")
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			os.MkdirAll(outputPath, 0755)
		}

		out_args := []string{}
		out_args = append(out_args, path)
		out_args = append(out_args, fmt.Sprintf("--go_out=%s", outputPath))
		out_args = append(out_args, "--go_opt=paths=source_relative")
		args := []string{}
		args = append(args, out_args...)
		args = append(args, search_args...)
		cmd := exec.Command("protoc", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(strings.Join(cmd.Args, " "))
			log.Println(string(out))
			log.Fatalf("Unable to run protoc on %s:%v", p.Name(), err)
		}
		genState.ProtoImportPath[p.Name()] = filepath.Join(outputPath, p.OutputPath())
	}

	for _, subPkg := range RosSubPackages {
		genPath := filepath.Join(subPkg.Dir(), subPkg.Name()) + ".go"
		cmd := exec.Command("gofmt", "-s", "-w", genPath)
		out, err := cmd.CombinedOutput()
		if len(out) > 0 {
			log.Print(string(out))
		}
		if err != nil {
			log.Println(strings.Join(cmd.Args, " "))
			log.Fatalf("unable to apply gofmt to %s: %v", subPkg.Name(), err)
		}
	}
	return
}

// a from -> to dependency relationship in an import graph
type edge struct {
	from string
	to   string
}

func deps(c *Context, pkgs []*ProtobufPackage) map[string][]string {
	ret := map[string][]string{}
	for _, p := range pkgs {
		pkg, ok := c.Universe[p.Path()]
		if !ok {
			log.Fatalf("unrecognized package %s", p.Path())
		}

		for _, d := range pkg.Imports {
			ret[p.Path()] = append(ret[p.Path()], d.Path)
		}
	}

	return ret
}

func removeEdgesFrom(node string, graph map[edge]struct{}) {
	for edge := range graph {
		if edge.from == node {
			delete(graph, edge)
		}
	}
}

func importOrder(deps map[string][]string) ([]string, error) {
	var remainingNodes = map[string]struct{}{}
	var graph = map[edge]struct{}{}
	for to, froms := range deps {
		remainingNodes[to] = struct{}{}
		for _, from := range froms {
			remainingNodes[from] = struct{}{}
			graph[edge{from: from, to: to}] = struct{}{}
		}
	}

	sorted := findAndRemoveNodesWithoutDependencies(remainingNodes, graph)
	for i := 0; i < len(sorted); i++ {
		node := sorted[i]
		removeEdgesFrom(node, graph)
		sorted = append(sorted, findAndRemoveNodesWithoutDependencies(remainingNodes, graph)...)
	}
	if len(remainingNodes) > 0 {
		return nil, fmt.Errorf("cycle: remaining nodes: %#v, remaining edges: %#v", remainingNodes, graph)
	}
	return sorted, nil
}

func findAndRemoveNodesWithoutDependencies(nodes map[string]struct{}, graph map[edge]struct{}) []string {
	roots := []string{}

	for node := range nodes {
		incoming := false

		for edge := range graph {
			if edge.to == node {
				incoming = true
				break
			}
		}

		if !incoming {
			delete(nodes, node)
			roots = append(roots, node)
		}
	}
	sort.Strings(roots)
	return roots
}

type positionOrder struct {
	pos      map[string]int
	elements []*ProtobufPackage
}

func (o positionOrder) Len() int {
	return len(o.elements)
}

func (o positionOrder) Less(i, j int) bool {
	return o.pos[o.elements[i].Path()] < o.pos[o.elements[j].Path()]
}

func (o positionOrder) Swap(i, j int) {
	o.elements[i], o.elements[j] = o.elements[j], o.elements[i]
}
