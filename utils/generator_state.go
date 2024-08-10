package utils

import (
	"go_agent/iface"
)

type GeneratorState struct {
	//map of package name to its Messages: pkg -> [msgName -> MessageDefinition]
	RosMsgPkgs map[string]map[string]interface{}
	// map of ros messages and their proto paths
	ProtoPkgPath map[string]string
	// map of ros messages and their generared proto paths
	ProtoImportPath map[string]string
}

var builders *builderUniverse

type builderFinder struct {
	GetBuilderFromName func(string) (iface.Builder, error)
	Generated          bool
}

func (b builderFinder) IsGenerated() bool {
	return b.Generated
}

type builderUniverse struct {
	builders map[string]builderFinder
}

func NewBuilder(name string, finder func(string) (iface.Builder, error)) builderFinder {

	if builders == nil {
		builders = &builderUniverse{}
		builders.builders = map[string]builderFinder{}
	}

	if _, ok := builders.builders[name]; ok {
		if builders.builders[name].GetBuilderFromName == nil && finder != nil {
			builders.builders[name] = builderFinder{GetBuilderFromName: finder, Generated: true}
		}
	} else {
		if finder != nil {
			builders.builders[name] = builderFinder{GetBuilderFromName: finder, Generated: true}
		} else {
			builders.builders[name] = builderFinder{Generated: false}
		}
	}

	return builders.builders[name]

}
