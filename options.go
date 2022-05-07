package goconcat

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

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

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) SetJSONOptions(jsonFilePath string) error {
	file, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return errors.WithStack(err)
	}

	err = jsoniter.UnmarshalFromString(string(file), &o)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o *Options) SetOptions(
	rootPath string,
	ignoredDirectories []Directory,
	filePrefix []PrefixType,
	destination string,
	deleteOldFiles bool,
	concatPkg bool,
	mockeryDestination bool,
	fileType []FileType,
) {
	o.RootPath = rootPath
	o.IgnoredDirectories = ignoredDirectories
	o.FilePrefix = filePrefix
	o.Destination = destination
	o.DeleteOldFiles = deleteOldFiles
	o.ConcatPackages = concatPkg
	o.MockeryDestination = mockeryDestination
	o.FileType = fileType
}
