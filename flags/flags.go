package flags

type flag struct {
	set   bool
	value string
}

func (f *flag) Set(value string) error {
	f.value = value
	f.set = true
	return nil
}

func (f *flag) String() string {
	return f.value
}

var (
	PathFlag flag
	Help     flag
)
