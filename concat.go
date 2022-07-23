package goconcat

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func GoConcat(options *Options) error {
	err := validateOptions(options)
	if err != nil {
		log.Println(err)
		return errors.WithStack(err)
	}

	filePaths, err := GetFilePaths(options)
	if err != nil {
		log.Println(err)
		return errors.WithStack(err)
	}

	filesToConcat, fileSet, err := ParseASTFiles(filePaths)
	if err != nil {
		log.Println(err)
		return errors.WithStack(err)
	}

	filesToSort, err := getFilesToSort(filesToConcat, options, fileSet)
	if err != nil {
		log.Println(err)
		return errors.WithStack(err)
	}

	for _, file := range filesToSort {
		originialPackage := file.Name.Name
		destinationPath := anyToString(options.Destination)

		if destinationPath == rootDirectory {
			destinationPath = goconcat
		}

		var pkgName string

		if splitStrings := strings.Split(destinationPath, "/"); len(splitStrings) > 0 {
			pkgName = splitStrings[len(splitStrings)-1]
		} else {
			pkgName = destinationPath
		}

		file.Name.Name = pkgName

		var buf bytes.Buffer
		if err := format.Node(&buf, fileSet, file); err != nil {
			log.Println(err)
			return errors.WithStack(err)
		}

		isValid := destinationDirIsValid(options.RootPath, destinationPath)

		if !isValid {
			if err := os.Mkdir(destinationPath, os.ModePerm); err != nil {
				log.Println(err)
				return errors.WithStack(err)
			}
		}

		var fileName string

		if !options.SplitFilesByPackage {
			fileName = goconcat
		} else {
			fileName = originialPackage
		}

		finalPath := getDestinationPath(destinationPath, file.Name.Name, FileGo, options, filePaths, fileName)

		if err := ioutil.WriteFile(finalPath, buf.Bytes(), os.ModePerm); err != nil {
			log.Println(err)
			return errors.WithStack(err)
		}
	}

	if options.DeleteOldFiles {
		if err := DeleteFiles(filePaths); err != nil {
			log.Println(err)
			return errors.WithStack(err)
		}
	}

	return nil
}

func getDestinationPath(
	destination string,
	packageName string,
	fileType FileType,
	options *Options,
	filePaths []string,
	fileName string,
) string {
	file := anyToString(fileType)

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

		return splitPath[0] + packageName + "/" + fmt.Sprintf("mocks_%s", packageName) + file
	}

	if options.SplitFilesByPackage {
		return "./" + destination + "/" + fileName + file
	}

	return "./" + destination + "/" + packageName + file
}

func validateOptions(options *Options) error {
	if options.RootPath == "" {
		return errors.WithStack(errNoRootPath)
	}

	if options.FilePrefix != nil && len(options.FilePrefix) < 1 {
		return errors.WithStack(errNoPrefix)
	}

	return nil
}

func concatImports(targetFile *ast.File, fileSet *token.FileSet, importStrings []string) {
	existingImports := make(map[string]string)

	for _, v := range targetFile.Imports {
		existingImports[v.Path.Value] = v.Path.Value
	}

	for _, importString := range importStrings {
		// skip import if it already exists
		if _, ok := existingImports[importString]; ok {
			continue
		}
		addImportToTargetFile(targetFile, importString)
	}
	ast.SortImports(fileSet, targetFile)
}

func addImportToTargetFile(targetFile *ast.File, target string) {
	for _, decl := range targetFile.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)

			if genDecl.Tok == token.IMPORT {
				spec := &ast.ImportSpec{
					Path: &ast.BasicLit{
						Value: target,
					},
				}
				genDecl.Specs = append(genDecl.Specs, spec)
			}
		}
	}
}

func getSpecsAndIndices(file *ast.File, tok token.Token) (specs []ast.Spec, declIndex []int) {
	for index, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl.Tok == tok {
				declIndex = append(declIndex, index)
				specs = append(specs, genDecl.Specs...)
			}
		}
	}

	return specs, declIndex
}

func concatSpecs(file *ast.File, specs []ast.Spec, tok token.Token) {
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl.Tok == tok {
				genDecl.Specs = append(genDecl.Specs, specs...)
			}
		}
	}
}

// removes decl from ast file
func removeDecl(file *ast.File, indices []int) {
	// get the first index of decl as the base decl
	indices = removeFromSlice(indices, 0)
	file.Decls = returnAllButIndices(file.Decls, indices)
}

func getFuncDeclFromFiles(files []*ast.File) []ast.Decl {
	var funcs []ast.Decl

	for _, file := range files {
		for _, decl := range file.Decls {
			switch decl.(type) {
			case *ast.FuncDecl:
				funcDecl := decl.(*ast.FuncDecl)
				funcs = append(funcs, funcDecl)
			}
		}
	}

	return funcs
}
