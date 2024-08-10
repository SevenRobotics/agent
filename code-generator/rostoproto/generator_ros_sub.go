package rostoproto

import (
	"io"
)

const (
	RosSubFileType = "go"
)

type RosSubGenerator struct {
	OutputFilename string

	OptionalBody []byte
}

func (rg RosSubGenerator) Name() string                                  { return rg.OutputFilename }
func (rg RosSubGenerator) Namers(*Context) NameSystems                   { return nil }
func (rg RosSubGenerator) Imports(*Context) []string                     { return []string{} }
func (rg RosSubGenerator) PackageVars(*Context) []string                 { return []string{} }
func (rg RosSubGenerator) GenerateType(*Context, *Type, io.Writer) error { return nil }
func (rg RosSubGenerator) Filename() string                              { return rg.OutputFilename + ".go" }
func (rg RosSubGenerator) FileType() string                              { return "go" }
func (rg RosSubGenerator) Finalize(*Context, io.Writer) error            { return nil }

func (rg RosSubGenerator) Init(c *Context, w io.Writer) error {
	_, err := w.Write(rg.OptionalBody)
	return err
}

var (
	_ = Generator(RosSubGenerator{})
)
