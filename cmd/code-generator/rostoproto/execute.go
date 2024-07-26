package rostoproto

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type DefaultFileType struct {
	Assemble func(io.Writer, *File)
}

func (ft DefaultFileType) AssembleFile(f *File, pathname string) error {
	log.Printf("Assembling file %q", pathname)

	destFile, err := os.Create(pathname)

	if err != nil {
		return err
	}

	defer destFile.Close()

	b := &bytes.Buffer{}
	ft.Assemble(b, f)

	_, err = destFile.Write(b.Bytes())
	return err
}

func (c *Context) ExecuteTargets(targets []Target) error {
	log.Printf("ExecuteTargets: %d targets", len(targets))

	var errs []error
	for _, tgt := range targets {
		if err := c.ExecuteTarget(tgt); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Some targets had errors : %w", errors.Join(errs...))
	}

	return nil
}

func (c *Context) addNameSystems(namers NameSystems) *Context {
	if namers == nil {
		return c
	}

	c2 := *c
	c2.Namers = NameSystems{}

	//copy existing namers
	for k, v := range c.Namers {
		c2.Namers[k] = v
	}

	//add new namers
	for name, namer := range namers {
		c2.Namers[name] = namer
	}

	return &c2

}

func (c *Context) executeBody(w io.Writer, generator Generator, t *Type) error {
	if err := generator.Init(c, w); err != nil {
		return err
	}
	if err := generator.GenerateType(c, t, w); err != nil {
		return err
	}
	if err := generator.Finalize(c, w); err != nil {
		return err
	}
	return nil
}

func (c *Context) ExecuteTarget(target Target) error {
	tgtDir := target.Dir()

	if tgtDir == "" {
		return fmt.Errorf("No target directory set for %v", target.Path())
	}

	if err := os.MkdirAll(tgtDir, 0755); err != nil {
		return err
	}

	files := map[string]*File{}
	for _, g := range target.Generators(c) {
		filetype := g.FileType()
		if len(filetype) == 0 {
			return fmt.Errorf("generator %q must specify a file type", g.Name())
		}

		genContext := c.addNameSystems(g.Namers(c))

		f := files[g.Filename()]

		if f == nil {
			f = &File{
				Name:        g.Filename(),
				FileType:    g.FileType(),
				PackageName: target.Name(),
				PackagePath: target.Path(),
				PackageDir:  target.Dir(),
				Header:      target.Header(g.Filename()),
				Imports:     map[string]struct{}{},
			}

			files[f.Name] = f
		} else if f.FileType != g.FileType() {
			return fmt.Errorf("file %q already has type %q, but generator %q wants to set type %q",
				f.Name, f.FileType, g.Name(), g.FileType())
		}

		if vars := g.PackageVars(genContext); len(vars) > 0 {
			for _, v := range vars {
				if _, err := fmt.Fprintf(&f.Vars, "%s\n", v); err != nil {
					return err
				}
			}
		}

		pack := c.Universe[target.Path()]
		msgMap := pack.MessageDefs[target.Path()]
		if msgDef, ok := msgMap[target.Name()]; ok {
			log.Printf("Message name: %v", msgDef.Name.Name)
			for k := range msgDef.Imports {
				log.Printf("Imports %v", k)
			}
			for i, field := range msgDef.Fields {
				log.Printf("Field %d : %v", i, field)
			}
			log.Printf("")
			if err := genContext.executeBody(&f.Body, g, &msgDef.Type); err != nil {
				return err
			}
		}
	}

	var errs []error
	if len(files) == 0 {
		log.Printf("Could not generate any files, perhaps no generator is set for this type")
	}
	for _, f := range files {
		finalPath := filepath.Join(tgtDir, f.Name)
		assembler, ok := c.FileTypes[f.FileType]
		if !ok {
			return fmt.Errorf("The filetype %q registered for file %q does not exit", f.FileType, f.Name)
		}

		if err := assembler.AssembleFile(f, finalPath); err != nil {
			errs = append(errs, err)
		}

	}

	if len(errs) > 0 {
		return fmt.Errorf("errors in target %q: %w", target, errors.Join(errs...))
	}

	return nil

}
