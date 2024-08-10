package rostoproto

import (
	"fmt"
	"io"
	"runtime"
	"text/template"
)

type SnippetWriter struct {
	w       io.Writer
	context *Context
	//left and right delimiters
	left, right string
	funcMap     template.FuncMap
	err         error
}

// c is used to make a function for every naming system, to which you can pass a type and get the corresponding name
func NewSnippetWriter(w io.Writer, c *Context, left, right string) *SnippetWriter {
	sw := &SnippetWriter{
		w:       w,
		context: c,
		left:    left,
		right:   right,
		funcMap: template.FuncMap{},
	}

	for name, namer := range c.Namers {
		sw.funcMap[name] = namer.Name
	}
	return sw
}

func (sw *SnippetWriter) Do(format string, args interface{}) *SnippetWriter {
	if sw.err != nil {
		return sw
	}

	_, file, line, _ := runtime.Caller(1)
	tmpl, err := template.
		New(fmt.Sprintf("%s:%d", file, line)).
		Delims(sw.left, sw.right).
		Funcs(sw.funcMap).
		Parse(format)
	if err != nil {
		sw.err = err
		return sw
	}
	err = tmpl.Execute(sw.w, args)
	if err != nil {
		sw.err = err
	}
	return sw
}

func (sw *SnippetWriter) Out() io.Writer {
	return sw.w
}
