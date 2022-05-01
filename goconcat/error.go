package goconcat

import "errors"

var (
	ErrReadingDirectories = errors.New("error: reading the directories")
	ErrNotEnoughFiles     = errors.New("error: not enough files were able to be read")
)
