package goconcat

import "errors"

var (
	errNoFilesDetected      = errors.New("error: there were no files detected")
	errNotEnoughtFiles      = errors.New("error: not enough files were found to concatenate")
	errReadingDirectories   = errors.New("error: reading the directories")
	errNotEnoughFiles       = errors.New("error: not enough files were able to be read")
	errBoolCouldNotBeParsed = errors.New("error: please make sure you pass in true or false")
	errNoFilePathForJson    = errors.New("error: you need to enter in a file path for your json file")
	errNoFilePath           = errors.New("error: no file path")
	errNoRootPath           = errors.New("error: a rootpath needs to be specified in the options")
	errNoPrefix             = errors.New("error: no prefix was specified in the options, please add one")
)
