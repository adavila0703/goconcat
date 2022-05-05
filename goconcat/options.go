package goconcat

type Options struct {
	RootPath           string       `json:"rootPath"`
	IgnoredDirectories []Directory  `json:"ignoredDirectories"`
	FilePrefix         []PrefixType `json:"filePrefix"`
	Destination        string       `json:"destination"`
	DeleteOldFiles     bool         `json:"deleteOldFiles"`
	ConcatPackages     bool         `json:"concatPkg"`
	MockeryDestination bool
	FileType           []FileType
}

func NewOptions(
	rootPath string,
	ignoredDirectories []Directory,
	filePrefix []PrefixType,
	destination string,
	deleteOldFiles bool,
	concatPkg bool,
	mockeryDestination bool,
	fileType []FileType,
) *Options {
	return &Options{
		RootPath:           rootPath,
		IgnoredDirectories: ignoredDirectories,
		FilePrefix:         filePrefix,
		Destination:        destination,
		DeleteOldFiles:     deleteOldFiles,
		ConcatPackages:     concatPkg,
		MockeryDestination: mockeryDestination,
		FileType:           fileType,
	}
}
