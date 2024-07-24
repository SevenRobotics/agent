package rostoproto

// The package name, path, and dir are required to be non-empty.
type SimpleTarget struct {
	// PkgName is the name of the resulting package (as in "package xxxx").
	// Required.
	PkgName string
	// PkgPath is the canonical Ros import-path of the resulting package (as in
	// "std_msgs/Header"). Required.
	PkgPath string
	// PkgDir is the location of the resulting package on disk (which may not
	// exist yet). It may be absolute or relative to CWD. Required.
	PkgDir string

	// HeaderComment is emitted at the top of every output file. Optional.
	HeaderComment []byte

	// GeneratorsFunc will be called to implement Target.Generators. Optional.
	GeneratorsFunc func(*Context) []Generator
}

func (st SimpleTarget) Name() string { return st.PkgName }
func (st SimpleTarget) Path() string { return st.PkgPath }
func (st SimpleTarget) Dir() string  { return st.PkgDir }
func (st SimpleTarget) Header(filename string) []byte {
	return st.HeaderComment
}

func (st SimpleTarget) Generators(c *Context) []Generator {
	if st.GeneratorsFunc != nil {
		return st.GeneratorsFunc(c)
	}
	return nil
}

var (
	_ = Target(SimpleTarget{})
)
