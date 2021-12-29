package goconcat

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrReadingDirectories = errors.New("error reading the directories")
)

func Goconcat() error {
	filePaths, err := getFilePaths(".", []string{".git"}, ".go", "mock_")
	if err != nil {
		return err
	}

	src, err := getSource(filePaths)
	if err != nil {
		fmt.Println(src)
		return err
	}

	test := `
	package main

	import (
		
		"fmt"
		"fmt"
	)
	
	`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", test, 0)
	if err != nil {
		panic(err)
	}

	// concat all imports
	if len(f.Imports) > 1 {
		f.Decls = concatImports(f)
	}

	ast.Print(fset, f)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())

	// ioutil.WriteFile("test.go", test, fs.ModeAppend)
	return nil
}

func getLengthOfVars(file *ast.File) int {
	var count int
	for _, decl := range file.Decls {
		switch v := decl.(type) {
		case *ast.GenDecl:
			if v.Tok == token.Token(token.VAR) {
				count++
			}
		}
	}
	return count
}

func getSource(paths []string) (string, error) {
	var src string

	for _, path := range paths {
		fileContents, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}

		src += string(fileContents)
	}

	return src, nil
}

func concatImports(file *ast.File) []ast.Decl {
	var newDecals []ast.Decl

	for index, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			if decl.(*ast.GenDecl).Tok == token.IMPORT && index != 0 {
				continue
			}

			if decl.(*ast.GenDecl).Tok == token.IMPORT && index == 0 {
				file.Decls[index].(*ast.GenDecl).Specs = concatSpec(file)
			}

			for _, value := range decl.(*ast.GenDecl).Specs {
				fmt.Println(value.(*ast.ImportSpec).Path.ValuePos)
			}
		}

		newDecals = append(newDecals, decl)
	}

	return newDecals
}

func concatSpec(file *ast.File) []ast.Spec {
	var newImports []ast.Spec

	for _, value := range file.Decls {
		switch v := value.(type) {
		case *ast.GenDecl:
			if v.Tok == token.IMPORT {
				newImports = append(newImports, v.Specs...)
			}
		}
	}

	return newImports
}

func getFilePaths(path string, ignoredDirectories []string, fileType string, prefix string) ([]string, error) {
	var filePaths []string

	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if checkDirectoryIgnore(path, ignoredDirectories) {
			return nil
		}

		if !info.IsDir() {
			if strings.HasSuffix(info.Name(), fileType) && strings.HasPrefix(info.Name(), prefix) {
				filePaths = append(filePaths, path)
			}
		}
		return nil
	})

	return filePaths, nil
}

func checkDirectoryIgnore(path string, ignoredDirectories []string) bool {
	for _, ignoreDirectory := range ignoredDirectories {
		if strings.Contains(path, ignoreDirectory) {
			return true
		}
	}
	return false
}

func fileToLines(filePath string) ([]string, error) {
	var lines []string

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}
