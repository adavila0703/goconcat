package goconcat

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"mockconcat/utils"

	"github.com/pkg/errors"
)

func GoConcat(
	path string,
	ignoredDirectories []utils.Directory,
	fileTypes []utils.FileType,
	prefix []utils.PrefixType,
) error {
	filePaths, err := GetFilePaths(
		path,
		ignoredDirectories,
		fileTypes,
		prefix,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	var files []*ast.File
	fileSet := token.NewFileSet()

	for _, path := range filePaths {
		fileContents, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}

		astFiles, err := parser.ParseFile(fileSet, "", fileContents, 0)
		if err != nil {
			return errors.WithStack(err)
		}

		files = append(files, astFiles)
	}

	concatFiles, err := ConcatFiles(files, fileSet)
	if err != nil {
		return errors.WithStack(err)
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, concatFiles); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
