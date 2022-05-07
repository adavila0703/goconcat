package goconcat

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoConcat(t *testing.T) {
	assert := assert.New(t)

	mockFileContent := map[string]string{
		"test_file_one.go": `package test
		var foo string
		`,
		"test_file_two.go": `package test
		var bar string
		`,
	}
	fileSet := token.NewFileSet()

	for key, value := range mockFileContent {
		file, err := parser.ParseFile(fileSet, key, value, 0)
		assert.NoError(err)

		var buf bytes.Buffer
		err = format.Node(&buf, fileSet, file)
		assert.NoError(err)

		err = ioutil.WriteFile(key, buf.Bytes(), os.ModePerm)
		assert.NoError(err)
	}

	options := NewOptions()
	options.SetOptions(
		".",
		nil,
		[]PrefixType{"test_"},
		".",
		true,
		false,
		false,
		[]FileType{FileGo},
	)

	err := GoConcat(options)
	assert.NoError(err)

	mockFilePaths := []string{"test.go"}

	files, parsedFileSet, err := ParseASTFiles(mockFilePaths)
	assert.NoError(err)

	var buf bytes.Buffer
	err = format.Node(&buf, parsedFileSet, files[0])
	assert.NoError(err)

	DeleteFiles(mockFilePaths)

	expectedResult := "package test\n\nvar (\n\tfoo string\n\tbar string\n)\n"

	assert.Equal(expectedResult, buf.String())
}

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

func TestRemoveDecl(t *testing.T) {
	assert := assert.New(t)
	var decls []ast.Decl

	funcs := []*ast.FuncDecl{
		{
			Name: &ast.Ident{
				Name: "test1",
			},
		},
		{
			Name: &ast.Ident{
				Name: "test2",
			},
		},
	}

	for _, f := range funcs {
		decls = append(decls, f)
	}

	file := &ast.File{
		Decls: decls,
	}

	removeDecl(file, []int{0, 1})

	var funcNames []string
	for _, decl := range file.Decls {
		switch f := decl.(type) {
		case *ast.FuncDecl:
			funcNames = append(funcNames, f.Name.Name)
		}
	}

	assert.Equal([]string{"test1"}, funcNames)
}

func TestGetFuncDeclFromFiles(t *testing.T) {
	assert := assert.New(t)
	var mockDecls []ast.Decl

	mockFuncs := []*ast.FuncDecl{
		{
			Name: &ast.Ident{
				Name: "test1",
			},
		},
		{
			Name: &ast.Ident{
				Name: "test2",
			},
		},
	}

	for _, f := range mockFuncs {
		mockDecls = append(mockDecls, f)
	}

	mockFile := []*ast.File{
		{
			Decls: mockDecls,
		},
	}

	decl := getFuncDeclFromFiles(mockFile)

	var funcNames []string
	for _, decl := range decl {
		switch f := decl.(type) {
		case *ast.FuncDecl:
			funcNames = append(funcNames, f.Name.Name)
		}
	}

	assert.Equal([]string{"test1", "test2"}, funcNames)
}
