package rostoproto

import (
	"bytes"
	"io"
)

// Target is a package into which code will be generated.
// A single target can have many generators, each of which emits one file.
// For instance, in our use case of parsing ros packages and converting to proto
// we can have generators for .msg, .srv and .action files.
type Target interface {
	//Name returns the package name in short (eg: visualization_msgs or uavcan_ros_bridge)
	Name() string

	//Path returns the package import path as in how to import this package
	Path() string

	//location of the package on the disk
	Dir() string

	//returns a header for the file. Only used for static comment markers
	Header(filename string) []byte

	Generators(*Context) []Generator
}

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
	Namers    NameSystems // used for naming packages & files
	Universe  Universe
	Inputs    []string //user specified packages
	FileTypes map[string]FileType
	parser    *Parser
}

// Generator is the contract for anything that wants to do auto-generation.
// The call order for the functions that take a Context is:
// 1. Filter()        // Subsequent calls see only types that pass this.
// 2. Namers()        // Subsequent calls see the namers provided by this.
// 3. PackageVars()
// 4. PackageConsts()
// 5. Init()
// 6. GenerateType()  // Called N times, once per type in the context's Order.
// 7. Imports()
//
// You may have multiple generators for the same file.
type Generator interface {
	// The name of this generator
	Name() string

	//Namers for this generator
	Namers(*Context) NameSystems

	Init(*Context, io.Writer) error

	Finalize(*Context, io.Writer) error

	PackageVars(*Context) []string

	GenerateType(*Context, *Type, io.Writer) error

	Imports(*Context) []string

	Filename() string

	FileType() string
}

func NewContext(p *Parser, nameSystem NameSystems) (*Context, error) {
	universe, err := p.NewUniverse()
	if err != nil {
		return nil, err
	}

	c := &Context{
		Namers:    NameSystems{},
		Universe:  universe,
		Inputs:    p.userRequestedPackages(),
		FileTypes: map[string]FileType{},
		parser:    p,
	}

	for name, systemNamer := range nameSystem {
		c.Namers[name] = systemNamer
	}

	return c, nil
}

func (c *Context) LoadPackages(patterns ...string) ([]*Package, error) {
	return c.parser.LoadPackagesTo(&c.Universe, patterns...)
}
