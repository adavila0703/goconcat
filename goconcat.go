package goconcat

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func GoConcat(
	options *Options,
) error {
	err := validateOptions(options)
	if err != nil {
		return errors.WithStack(err)
	}

	filePaths, err := GetFilePaths(
		options.RootPath,
		options.IgnoredDirectories,
		options.FileType,
		options.FilePrefix,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	var filesToConcat []*ast.File
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

		filesToConcat = append(filesToConcat, astFiles)
	}

	filesToSort, err := GetFilesToSort(filesToConcat, options, fileSet)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range filesToSort {
		var buf bytes.Buffer
		if err := format.Node(&buf, fileSet, file); err != nil {
			return errors.WithStack(err)
		}

		des := AnyToString(options.Destination)
		isValid := DestinationDirIsValid(options.RootPath, des)

		if !isValid && !options.MockeryDestination {
			if err := os.Mkdir(des, os.ModePerm); err != nil {
				return errors.WithStack(err)
			}
		}

		finalPath := GetDestinationPath(des, file.Name.Name, FileGo, options, filePaths)

		if err := ioutil.WriteFile(finalPath, buf.Bytes(), os.ModePerm); err != nil {
			return errors.WithStack(err)
		}
	}

	if options.DeleteOldFiles {
		if err := DeleteFiles(filePaths); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func GetFilesToSort(files []*ast.File, options *Options, fileSet *token.FileSet) ([]*ast.File, error) {
	var filesToSort []*ast.File

	if options.ConcatPackages {
		filePackageMap := make(map[string][]*ast.File)

		for _, file := range files {
			packageName := file.Name.Name

			if _, ok := filePackageMap[packageName]; !ok {
				filePackageMap[packageName] = []*ast.File{file}
			} else {
				filePackageMap[packageName] = append(filePackageMap[packageName], file)
			}
		}

		for _, files := range filePackageMap {
			concatFiles, err := ConcatFiles(files, fileSet)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			filesToSort = append(filesToSort, concatFiles)
		}

	} else {
		concatFiles, err := ConcatFiles(files, fileSet)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		filesToSort = append(filesToSort, concatFiles)
	}

	return filesToSort, nil
}

func GetDestinationPath(
	destination string,
	packageName string,
	fileType FileType,
	options *Options,
	filePaths []string,
) string {
	file := AnyToString(fileType)

	if options.MockeryDestination {
		findPackage := regexp.MustCompile(packageName)

		var splitPath []string
		for _, path := range filePaths {
			if findPackage.Match([]byte(path)) {
				splitPath = strings.Split(path, packageName)
			}
		}

		if splitPath[0] == "" {
			return packageName + "/" + fmt.Sprintf("mocks_%s", packageName) + file
		}

		return splitPath[0] + "/" + fmt.Sprintf("mocks_%s", packageName) + file
	}

	return "./" + destination + "/" + packageName + file
}

func validateOptions(options *Options) error {
	if options.FileType == nil {
		options.FileType = []FileType{FileGo}
	}

	return nil
}
