package flags

import "flag"

var (
	Test    = flag.String("something", "woop", "anything")
	Pkgname = flag.String("p", "", "package name to use in output file")
	Prefix  = flag.String("x", "", "prefix to add to all top-level names")
	Notest  = flag.Bool("n", false, "ignore test files")
	Kill    = flag.Bool("k", false, "delete concatenated files from disk")
)
