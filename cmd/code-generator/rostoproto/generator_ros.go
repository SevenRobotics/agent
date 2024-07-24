package rostoproto

import (
	"io"
)

const (
	RosMsgFileType = "msg"
)

type RosMsgGenerator struct {
	OutputFilename string

	OptionalBody []byte
}

func (rg RosMsgGenerator) Name() string                                  { return rg.OutputFilename }
func (rg RosMsgGenerator) Namers(*Context) NameSystems                   { return nil }
func (rg RosMsgGenerator) Imports(*Context) []string                     { return []string{} }
func (rg RosMsgGenerator) PackageVars(*Context) []string                 { return []string{} }
func (rg RosMsgGenerator) GenerateType(*Context, *Type, io.Writer) error { return nil }
func (rg RosMsgGenerator) Filename() string                              { return rg.OutputFilename + ".proto" }
func (rg RosMsgGenerator) FileType() string                              { return RosMsgFileType }
func (rg RosMsgGenerator) Finalize(*Context, io.Writer) error            { return nil }

func (rg RosMsgGenerator) Init(c *Context, w io.Writer) error {
	_, err := w.Write(rg.OptionalBody)
	return err
}

var (
	_ = Generator(RosMsgGenerator{})
)
