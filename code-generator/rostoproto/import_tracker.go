package rostoproto

import "sort"

type DefaultImportTracker struct {
	pathToName  map[string]string
	nameToPath  map[string]string
	local       Name
	LocalName   func(Name) string
	PrintImport func(string, string) string
}

func NewDefaultImportTracker(local Name) DefaultImportTracker {
	return DefaultImportTracker{
		pathToName: map[string]string{},
		nameToPath: map[string]string{},
		local:      local,
	}
}

func (tracker *DefaultImportTracker) AddTypes(types ...*Type) {
	for _, t := range types {
		tracker.AddType(t)
	}
}

func (tracker *DefaultImportTracker) AddType(t *Type) {
	tracker.AddSymbol(t.Name)
}

func (tracker *DefaultImportTracker) AddSymbol(name Name) {
	if len(name.Package) == 0 {
		return
	}

	path := name.Path

	if len(path) == 0 {
		path = name.Package
	}

	if _, ok := tracker.pathToName[path]; ok {
		return
	}

	n := tracker.LocalName(name)
	tracker.nameToPath[n] = path
	tracker.pathToName[path] = n
}

func (tracker *DefaultImportTracker) ImportLines() []string {
	importPaths := []string{}

	for path := range tracker.pathToName {
		importPaths = append(importPaths, path)
	}

	sort.Strings(importPaths)

	out := []string{}
	for _, path := range importPaths {
		out = append(out, tracker.PrintImport(path, tracker.pathToName[path]))
	}
	return out
}

func (tracker *DefaultImportTracker) LocalNameOf(path string) string {
	return tracker.pathToName[path]
}

func (tracker *DefaultImportTracker) PathOf(localName string) (string, bool) {
	name, ok := tracker.nameToPath[localName]
	return name, ok
}
