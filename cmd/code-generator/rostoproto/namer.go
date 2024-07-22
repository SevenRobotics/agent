package rostoproto

import (
	"fmt"
	"strings"
)

func protoSafePackage(name string) string {
	pkg := strings.ReplaceAll(name, "/", ".")
	return strings.ReplaceAll(pkg, "-", "_")
}

func protoSafePackageMsg(name string) string {
	msg := strings.Split(name, ".")[0]
	return protoSafePackage(msg)
}

type localNamer struct {
	localPackage Name
}

func (n localNamer) Name(t *Type) string {
	if len(n.localPackage.Package) != 0 && n.localPackage.Package == t.Name.Package {
		return t.Name.Name
	}

	if strings.Contains(t.Name.Package, ".") {
		return fmt.Sprintf(".%s", t.Name)
	}

	return t.Name.String()
}

type protoBufNamer struct {
	packages []*ProtobufPackage
	// the key for this will be a ros msg import path
	packagesByPath map[string]*ProtobufPackage
}

func NewProtobufNamer() *protoBufNamer {
	return &protoBufNamer{
		packagesByPath: make(map[string]*ProtobufPackage),
	}
}

func (n *protoBufNamer) Name(t *Type) string {
	return t.Name.String()
}

func (n *protoBufNamer) Add(p *ProtobufPackage) {
	if _, ok := n.packagesByPath[p.Path()]; !ok {
		n.packagesByPath[p.Path()] = p
		n.packages = append(n.packages, p)
	}
}
