package utils

type Options struct {
	ConcatPackages     bool
	MockeryDestination bool
}

func NewOptions(concatPkg bool, mockeryDestination bool) *Options {
	return &Options{
		ConcatPackages:     concatPkg,
		MockeryDestination: mockeryDestination,
	}
}
