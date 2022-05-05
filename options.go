package goconcat

import "github.com/adavila0703/goconcat/internal/utils"

type Options struct {
	RootPath           string             `json:"rootPath"`
	IgnoredDirectories []utils.Directory  `json:"ignoredDirectories"`
	FilePrefix         []utils.PrefixType `json:"filePrefix"`
	Destination        string             `json:"destination"`
	DeleteOldFiles     bool               `json:"deleteOldFiles"`
	ConcatPackages     bool               `json:"concatPkg"`
	MockeryDestination bool
	FileType           []utils.FileType
}

func NewOptions(
	rootPath string,
	ignoredDirectories []utils.Directory,
	filePrefix []utils.PrefixType,
	destination string,
	deleteOldFiles bool,
	concatPkg bool,
	mockeryDestination bool,
	fileType []utils.FileType,
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
