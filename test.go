// main package.
package main

import (
	// "io"
	// "os"
	//
	// "github.com/alecthomas/kong"
	//
	"go_agent/cmd/code-generator/rostoproto"
)

var cli struct {
	GoPackage  string `name:"gopackage" help:"Go package name" default:"main"`
	RosPackage string `name:"rospackage" help:"ROS package name" default:"my_package"`
	Path       string `arg:"" help:"path pointing to a ROS message"`
}

// func run(args []string, output io.Writer) error {
// 	parser, err := kong.New(&cli,
// 		kong.Description("Convert ROS messages into Go structs."),
// 		kong.UsageOnError())
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = parser.Parse(args)
// 	if err != nil {
// 		return err
// 	}
//
// 	return message.ImportMessage(cli.Path, cli.GoPackage, cli.RosPackage, output)
// }

var g = rostoproto.NewGen()

func main() {
	rostoproto.Run(g)
	// // err := run(os.Args[1:], os.Stdout)
	// // if err != nil {
	// // 	fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
	// // 	os.Exit(1)
	// // }
	// packages, err := util.FindRosPackages()
	// if err != nil {
	// 	fmt.Println("Error finding ros packages")
	// }
	//
	// p := rostoproto.NewParser()
	// pkgs, err := p.LoadPackages(packages...)
	//
	// for _, pkg := range pkgs {
	// 	fmt.Printf("Package Name: %v Path: %v", pkg.Name, pkg.Dir)
	// 	for _, msgDef := range pkg.MessageDefs {
	// 		fmt.Printf("Message PackageName: %v\n", msgDef.RosPkgName)
	// 		fmt.Printf("Message Name: %v\n", msgDef.Name)
	// 		for i, field := range msgDef.Fields {
	// 			fmt.Printf("Field Definition %v : %v\n", i, field)
	// 		}
	//
	// 		for i, def := range msgDef.Definitions {
	// 			fmt.Printf("Message Definition %v: %v\n", i, def)
	// 		}
	//
	// 		fmt.Printf("Message DefinitionStr: %v\n", msgDef.DefinitionsStr)
	// 	}
	// }

}
