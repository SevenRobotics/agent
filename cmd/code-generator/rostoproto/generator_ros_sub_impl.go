package rostoproto

import (
	"fmt"
	"io"
	"strings"
)

type rosSubGen struct {
	RosSubGenerator
	localPackage    Name
	localRosPackage Name
	imports         ImportTracker
}

func (r *rosSubGen) PackageVars(c *Context) []string {
	return []string{
		fmt.Sprintf("package %s\n", r.localRosPackage),
	}
}

func (r *rosSubGen) FileName() string { return r.OutputFilename + ".go" }
func (r *rosSubGen) FileType() string { return "go" }

func (r *rosSubGen) Namers(c *Context) NameSystems {
	return NameSystems{
		"local": localNamer{r.localPackage},
	}
}

func (r *rosSubGen) Imports(c *Context) (imports []string) {
	lines := []string{}

	return lines
}

func (r *rosSubGen) GenerateType(c *Context, t *Type, w io.Writer) error {
	sw := NewSnippetWriter(w, c, "$", "$")
	b := subBodyGen{
		locator: &rosSubLocator{
			namer:           c.Namers["go"].(RosSubNamer),
			universe:        c.Universe,
			localRosPackage: r.localRosPackage.String(),
		},
		localPackage: r.localPackage,
		t:            t,
	}

	switch t.Kind {
	case RosMsgType:
		return b.doRosSub(sw)
	default:
		return b.unknown(sw)
	}

}

type subBodyGen struct {
	locator      RosSubLocator
	localPackage Name
	t            *Type
}

type RosSubLocator interface {
	RosTypeForName(name *Name) *Type
	CastTypeName(name Name) string
	GetMessageDef(name Name, msgName string) (*MessageDefinition, error)
	RosPackageName() string
}

type RosSubNamer interface {
	RosNameToGoName(name Name) Name
}

type rosSubLocator struct {
	namer           RosSubNamer
	tracker         ImportTracker
	universe        Universe
	localRosPackage string
}

func (r rosSubLocator) CastTypeName(name Name) string {
	if name.Package == r.localRosPackage {
		return name.Name
	}
	return name.String()
}

func (r rosSubLocator) RosPackageName() string {
	return r.localRosPackage
}

func (r rosSubLocator) GetMessageDef(name Name, msgName string) (*MessageDefinition, error) {
	if t, ok := r.universe[name.Path]; ok {
		if msg, ok := t.MessageDefs[name.Path]; ok {
			ret := msg[msgName]

			return ret, nil
		}
	}
	return nil, fmt.Errorf("Message type %v not found in package %v", msgName, name.Package)
}

func (r rosSubLocator) RosTypeForName(name *Name) *Type {
	if name == nil {
		return nil
	}

	if len(name.Package) == 0 {
		name.Package = r.localRosPackage
	}

	return r.universe.Type(*name)
}

func (b subBodyGen) doRosSub(sw *SnippetWriter) error {
	if len(b.t.Name.Name) == 0 {
		return nil
	}

	out := sw.Out()
	fmt.Fprintf(out, "\n\n")
	pkgName := strings.Split(b.locator.RosPackageName(), ".")
	b.t.Name.Package = pkgName[0]
	sw.Do(
		`import (
    "github.com/bluenviron/goroslib/v2"
    "go_agent/telemetry/gengo/ros"
    )
    
    var msgStream chan *$.Name.Package$.$.Name.Name$
    var initialized bool = false

    func msgCallback(msg *$.Name.Package$.$.Name.Name$) {
      if !initialized {
        return 
      }

      if msg != nil {
        msgStream <- msg 
      }
    }

    func Init(topicName string, n *goroslib.Node, msgS chan *$.Name.Package$.$.Name.Name$, done chan int) error {
      msgStream = msgS 
      sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
        Node: n,
        Topic: topicName,
        Callback: msgCallback, 
      })

      if err != nil {
        return fmt.Errorf("Failed to create subscriber for topic %q: %v",topicName,err)
      }
      initialized = true
      // wait for done 
      <-done
    }
    `, b.t)
	return nil
}

func (b subBodyGen) unknown(sw *SnippetWriter) error {
	return fmt.Errorf("cannot generator ros sub for %v", b.t)
}

func assembleRosSubFile(w io.Writer, file *File) {
	w.Write(file.Header)
	if len(file.PackageName) != 0 {
		fmt.Fprintf(w, "package %s\n", file.PackagePath)
	}
	w.Write(file.Body.Bytes())
}

func NewRosSubFile() *DefaultFileType {
	return &DefaultFileType{
		Assemble: assembleRosSubFile,
	}
}
