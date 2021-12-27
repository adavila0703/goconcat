package goconcat

import (
	"bufio"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
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
	const hello = "Hello, World!" // line comment 1
	
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

	fmt.Println(f.Imports)

	// ioutil.WriteFile("test.go", test, fs.ModeAppend)
	return nil
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
