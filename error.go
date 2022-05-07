package goconcat

import "errors"

var (
	ErrNoFilesDetected      = errors.New("error: there were no files detected")
	ErrNotEnoughtFiles      = errors.New("error: not enough files were found to concatenate")
	ErrReadingDirectories   = errors.New("error: reading the directories")
	ErrNotEnoughFiles       = errors.New("error: not enough files were able to be read")
	ErrBoolCouldNotBeParsed = errors.New("error: please make sure you pass in true or false")
	ErrNoFilePathForJson    = errors.New("error: you need to enter in a file path for your json file")
)
