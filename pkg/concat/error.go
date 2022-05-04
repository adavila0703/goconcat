package concat

import "errors"

var (
	ErrNoFilesDetected    = errors.New("error: there were no files detected")
	ErrNotEnoughtFiles    = errors.New("error: not enough files were found to concatenate")
	ErrReadingDirectories = errors.New("error: reading the directories")
	ErrNotEnoughFiles     = errors.New("error: not enough files were able to be read")
)
