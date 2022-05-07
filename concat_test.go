package goconcat

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDestinationPath(t *testing.T) {
	assert := assert.New(t)

	mockOptions := NewOptions()

	path := getDestinationPath("goconcat/test", "test", FileGo, mockOptions, []string{"goconcat/test/mock_file.go"})
	assert.Equal("./goconcat/test/test.go", path)
}
func TestGetDestinationPath_MockeryDestination(t *testing.T) {
	assert := assert.New(t)

	mockOptions := NewOptions()
	mockOptions.MockeryDestination = true

	path := getDestinationPath(".", "test", FileGo, mockOptions, []string{"goconcat/test/mock_file.go"})
	assert.Equal("goconcat/test/mocks_test.go", path)
}

func TestGetDestinationPath_MockeryDestination_NoPath(t *testing.T) {
	assert := assert.New(t)

	mockOptions := NewOptions()
	mockOptions.MockeryDestination = true

	path := getDestinationPath(".", "test", FileGo, mockOptions, []string{"test/mock_file.go"})
	assert.Equal("test/mocks_test.go", path)
}

func TestValidateOptions(t *testing.T) {
	assert := assert.New(t)
	options := NewOptions()

	err := validateOptions(options)
	assert.NoError(err)
	assert.Equal(FileGo, options.FileType[0])
}

func TestConcatImports(t *testing.T) {
	assert := assert.New(t)

	importStrings := []string{`"ast"`, `"fmt"`, `"fmt"`}
	fileSet := token.NewFileSet()

	mockFile := `package test
		import "fmt"
		`
	targetFile, err := parser.ParseFile(fileSet, "", mockFile, 0)
	assert.NoError(err)

	concatImports(targetFile, fileSet, importStrings)

	var specs []ast.Spec
	for _, decl := range targetFile.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)

			if genDecl.Tok == token.IMPORT {
				specs = append(specs, genDecl.Specs...)
			}
		}
	}

	var importValues []string
	for _, spec := range specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			importValues = append(importValues, s.Path.Value)
		}
	}

	assert.Equal([]string{"\"fmt\"", "\"ast\""}, importValues)
}

func TestAddImportToTargetFile(t *testing.T) {
	assert := assert.New(t)

	fileSet := token.NewFileSet()

	targetContents := `package test
		import "fmt"
		`
	targetFile, err := parser.ParseFile(fileSet, "", targetContents, 0)
	assert.NoError(err)

	addImportToTargetFile(targetFile, `"ast"`)

	var specs []ast.Spec
	for _, decl := range targetFile.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)

			if genDecl.Tok == token.IMPORT {
				specs = append(specs, genDecl.Specs...)
			}
		}
	}

	var importValues []string
	for _, spec := range specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			importValues = append(importValues, s.Path.Value)
		}
	}

	assert.Equal([]string{"\"fmt\"", "\"ast\""}, importValues)
}

func TestGetSpecsAndIndices(t *testing.T) {
	assert := assert.New(t)

	fileSet := token.NewFileSet()

	targetContents := `package test
		import (
			"fmt"
			"ast"
			"token"
		)
		`
	targetFile, err := parser.ParseFile(fileSet, "", targetContents, 0)
	assert.NoError(err)

	specs, declIndex := getSpecsAndIndices(targetFile, token.IMPORT)

	var importValues []string
	for _, spec := range specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			importValues = append(importValues, s.Path.Value)
		}
	}

	assert.Equal([]string{"\"fmt\"", "\"ast\"", "\"token\""}, importValues)
	assert.Equal([]int{0}, declIndex)
}

func TestConcatSpecs(t *testing.T) {
	assert := assert.New(t)

	targetContents := `package test
	import (
		"fmt"
	)
	`
	fileSet := token.NewFileSet()
	targetFile, err := parser.ParseFile(fileSet, "", targetContents, 0)
	assert.NoError(err)

	mockSpecs := []*ast.ImportSpec{
		{
			Path: &ast.BasicLit{
				Value: `"ast"`,
			},
		},
		{
			Path: &ast.BasicLit{
				Value: `"token"`,
			},
		},
	}

	var astSpec []ast.Spec

	for _, spec := range mockSpecs {
		astSpec = append(astSpec, spec)
	}

	concatSpecs(targetFile, astSpec, token.IMPORT)

	var importValues []string
	for _, spec := range astSpec {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			importValues = append(importValues, s.Path.Value)
		}
	}

	assert.Equal([]string{"\"ast\"", "\"token\""}, importValues)
}
