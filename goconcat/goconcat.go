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
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	ErrReadingDirectories = errors.New("error reading the directories")
)

func Goconcat() error {
	// filePaths, err := getFilePaths(".", []string{".git"}, ".go", "mock_")
	// if err != nil {
	// 	return err
	// }

	// test, err := ioutil.ReadFile(filePaths[0])
	// if err != nil {
	// 	return err
	// }

	// test2, err := ioutil.ReadFile(filePaths[1])
	// if err != nil {
	// 	return err
	// }

	// _, _ := fileToLines(filePaths[0])
	// test = append(test, test2...)

	src := `
	// This is the package comment.
	package main

	import (
		"fmt"
		"os"
	)

	// This comment is associated with the hello constant.
	const hello = "Hello, package!" // line comment 1
	
	// This comment is associated with the foo variable.
	var foo = hello // line comment 2
	
	// This comment is associated with the main function.
	func main() {
		fmt.Println(hello) // line comment 3
	}
	`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, parser.ImportsOnly)
	if err != nil {
		panic(err)
	}

	ast.Print(fset, f)

	var imports []string

	for _, value := range f.Imports {
		imports = append(imports, value.Path)
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())

	// test := ast.NewCommentMap(fset, f, f.Comments)

	// ioutil.WriteFile("test.go", test, fs.ModeAppend)
	return nil
}

func addImport(f *ast.File, ipath string) (added bool) {
	if imports(f, ipath) {
		return false
	}

	// Determine name of import.
	// Assume added imports follow convention of using last element.
	_, name := path.Split(ipath)

	// Rename any conflicting top-level references from name to name_.
	renameTop(f, name, name+"_")

	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(ipath),
		},
	}

	// Find an import decl to add to.
	var (
		bestMatch  = -1
		lastImport = -1
		impDecl    *ast.GenDecl
		impIndex   = -1
	)
	for i, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.IMPORT {
			lastImport = i

			// Compute longest shared prefix with imports in this block.
			for j, spec := range gen.Specs {
				impspec := spec.(*ast.ImportSpec)
				n := matchLen(importPath(impspec), ipath)
				if n > bestMatch {
					bestMatch = n
					impDecl = gen
					impIndex = j
				}
			}
		}
	}

	// If no import decl found, add one after the last import.
	if impDecl == nil {
		impDecl = &ast.GenDecl{
			Tok: token.IMPORT,
		}
		f.Decls = append(f.Decls, nil)
		copy(f.Decls[lastImport+2:], f.Decls[lastImport+1:])
		f.Decls[lastImport+1] = impDecl
	}

	// Ensure the import decl has parentheses, if needed.
	if len(impDecl.Specs) > 0 && !impDecl.Lparen.IsValid() {
		impDecl.Lparen = impDecl.Pos()
	}

	insertAt := impIndex + 1
	if insertAt == 0 {
		insertAt = len(impDecl.Specs)
	}
	impDecl.Specs = append(impDecl.Specs, nil)
	copy(impDecl.Specs[insertAt+1:], impDecl.Specs[insertAt:])
	impDecl.Specs[insertAt] = newImport
	if insertAt > 0 {
		// Assign same position as the previous import,
		// so that the sorter sees it as being in the same block.
		prev := impDecl.Specs[insertAt-1]
		newImport.Path.ValuePos = prev.Pos()
		newImport.EndPos = prev.Pos()
	}

	f.Imports = append(f.Imports, newImport)
	return true
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
