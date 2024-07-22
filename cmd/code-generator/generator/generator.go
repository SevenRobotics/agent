package generator

import (
	"bytes"
	"fmt"
)

type File struct {
	Name        string
	FileType    string
	PackageName string
	Header      []byte
	PackagePath string
	PackageDir  string
	Imports     map[string]struct{}
	Vars        bytes.Buffer
	Consts      bytes.Buffer
	Body        bytes.Buffer
}

type FileType interface {
	AssembleFile(f *File, path string) error
}

type Context struct {
	Namers    namer.NameSystems // used for naming packages & files
	Universe  namer.Universe
	Inputs    []string //user specified packages
	FileTypes map[string]FileType
	parser    *parser.Parser
}

type genProtoIDL struct {
	RosMsgGenerator
	localPackage    string
	localRosPackage string
	// imports         namer.ImportTracker

	generateAll    bool
	omitFieldTypes map[string]struct{}
}

func (g *genProtoIDL) PackageVars(c *Context) []string {
	return []string{
		fmt.Sprintf("option go_package=%q;", g.localRosPackage),
	}
}

func (g *genProtoIDL) FileName() string { return g.OutputFilename + ".proto" }
func (g *genProtoIDL) FileType() string { return "protoidl" }

// func (g* genProtoIDL) Namers(c* Context) return namer.NameSystems {
//   return namer.NameSystems {
//     "local": localNamer{g.localPackage},
//   }
// }
//
// func (g* genProtoIDL) Imports(c *Context) (imports []string) {
//   lines := []string{}
//
//   for _, line := range g.imports.ImportLines() {
//     lines = append(lines, line)
//   }
//
//   return lines
// }
