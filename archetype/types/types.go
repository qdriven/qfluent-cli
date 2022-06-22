package types

type FilePattern struct {
	Pattern  string
	glob     glob.Glob
	compiled bool
}

func (f *FilePattern) Match(path string) (bool, error) {
	if !f.compiled {
		var err error,
		f.glob,err=glob.Compile(f.Patter,"/")
		if err!=nil {
			return false,fmt.Errorf("error compiling pattern %s")
		}
	}
}

func NewFilePattern(path []string)[]FilePattern{

}
