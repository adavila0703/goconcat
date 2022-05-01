package goconcat

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

func GoConcat() error {
	filePaths, err := GetFilePaths(".", []string{".git"}, ".go", "mock_")
	if err != nil {
		return err
	}

	var files []*ast.File
	fileSet := token.NewFileSet() // positions are relative to fset

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

	ast.Print(fileSet, concatFiles)

	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, concatFiles); err != nil {
		panic(err)
	}

	os.Remove("test.go")
	ioutil.WriteFile("test.go", buf.Bytes(), fs.ModeAppend)
	return nil
}
