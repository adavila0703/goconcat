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

func TestConcatFiles_OneFile(t *testing.T) {
	assert := assert.New(t)

	mockFiles := []*ast.File{
		{
			Name: &ast.Ident{
				Name: "test",
			},
		},
	}

	mockFileSet := token.NewFileSet()

	file, err := ConcatFiles(mockFiles, mockFileSet, nil)

	assert.NoError(err)
	assert.Equal(mockFiles[0], file)
}

func TestConcatFiles_Imports(t *testing.T) {
	assert := assert.New(t)

	mockFiles := []string{`
		package test

		import (
			"fmt"
		)
		`,
		`
		package test

		import (
			"fmt"
			"go/token"
		)
		`,
	}

	expectedFileContents := `
	package test

	import (
		"fmt"
		"go/token"
	)
	`

	mockFileSet := token.NewFileSet()

	var files []*ast.File

	for _, mockFile := range mockFiles {
		file, err := parser.ParseFile(mockFileSet, "", mockFile, 0)
		assert.NoError(err)

		files = append(files, file)
	}

	concatFile, err := ConcatFiles(files, mockFileSet, nil)
	assert.NoError(err)

	var output bytes.Buffer
	if err := format.Node(&output, mockFileSet, concatFile); err != nil {
		assert.NoError(err)
	}

	expectedFile, err := parser.ParseFile(mockFileSet, "", expectedFileContents, 0)
	assert.NoError(err)

	var expectedOutput bytes.Buffer
	if err := format.Node(&expectedOutput, mockFileSet, expectedFile); err != nil {
		assert.NoError(err)
	}

	assert.Equal(expectedOutput.String(), output.String())
}

func TestConcatFiles_Vars(t *testing.T) {
	assert := assert.New(t)

	mockFiles := []string{`
		package test

		var (
			name = "name"
		)
		`,
		`
		package test

		var (
			age = 30
		)
		`,
	}

	expectedFileContents := `
	package test

	var (
		name = "name"
		age = 30
	)
	`

	mockFileSet := token.NewFileSet()

	var files []*ast.File

	for _, mockFile := range mockFiles {
		file, err := parser.ParseFile(mockFileSet, "", mockFile, 0)
		assert.NoError(err)

		files = append(files, file)
	}

	concatFile, err := ConcatFiles(files, mockFileSet, nil)
	assert.NoError(err)

	var output bytes.Buffer
	if err := format.Node(&output, mockFileSet, concatFile); err != nil {
		assert.NoError(err)
	}

	expectedFile, err := parser.ParseFile(mockFileSet, "", expectedFileContents, 0)
	assert.NoError(err)

	var expectedOutput bytes.Buffer
	if err := format.Node(&expectedOutput, mockFileSet, expectedFile); err != nil {
		assert.NoError(err)
	}

	assert.Equal(expectedOutput.String(), output.String())
}

func TestConcatFiles_Const(t *testing.T) {
	assert := assert.New(t)

	mockFiles := []string{`
		package test

		const (
			name = "name"
		)
		`,
		`
		package test

		const (
			age = 30
		)
		`,
	}

	expectedFileContents := `
	package test

	const (
		name = "name"
		age = 30
	)
	`

	mockFileSet := token.NewFileSet()

	var files []*ast.File

	for _, mockFile := range mockFiles {
		file, err := parser.ParseFile(mockFileSet, "", mockFile, 0)
		assert.NoError(err)

		files = append(files, file)
	}

	concatFile, err := ConcatFiles(files, mockFileSet, nil)
	assert.NoError(err)

	var outPut bytes.Buffer
	if err := format.Node(&outPut, mockFileSet, concatFile); err != nil {
		assert.NoError(err)
	}

	expectedFile, err := parser.ParseFile(mockFileSet, "", expectedFileContents, 0)
	assert.NoError(err)

	var expectedOutput bytes.Buffer
	if err := format.Node(&expectedOutput, mockFileSet, expectedFile); err != nil {
		assert.NoError(err)
	}

	assert.Equal(expectedOutput.String(), outPut.String())
}

func TestConcatFiles_Types(t *testing.T) {
	assert := assert.New(t)

	mockFiles := []string{`
		package test

		type IPerson interface {
			SayHello()
		}
	
		type Person struct {
			Name string
			Age  int
		}
		`,
		`
		package test

		type IDog interface {
			Bark()
		}
	
		type Dog struct {
			Name string
			Age  int
		}
		`,
	}

	expectedFileContents := `
	package test

	type (
		IPerson interface {
			SayHello()
		}
	
		Person struct {
			Name string
			Age  int
		}
		IDog interface {
			Bark()
		}
	
		Dog struct {
			Name string
			Age  int
		}
	)
	`

	mockFileSet := token.NewFileSet()

	var files []*ast.File

	for _, mockFile := range mockFiles {
		file, err := parser.ParseFile(mockFileSet, "", mockFile, 0)
		assert.NoError(err)

		files = append(files, file)
	}

	concatFile, err := ConcatFiles(files, mockFileSet, nil)
	assert.NoError(err)

	var outPut bytes.Buffer
	if err := format.Node(&outPut, mockFileSet, concatFile); err != nil {
		assert.NoError(err)
	}

	expectedFile, err := parser.ParseFile(mockFileSet, "", expectedFileContents, 0)
	assert.NoError(err)

	var expectedOutput bytes.Buffer
	if err := format.Node(&expectedOutput, mockFileSet, expectedFile); err != nil {
		assert.NoError(err)
	}

	assert.Equal(expectedOutput.String(), outPut.String())
}

func TestConcatFiles_Func(t *testing.T) {
	assert := assert.New(t)

	mockFiles := []string{`
		package test

		func (p *Person) SayHello() {
			fmt.Println("hello")
		}
		`,
		`
		package test

		func PrintToken() {
			tokens := []token.Token{
				token.VAR,
				token.CONST,
			}
		
			for _, token := range tokens {
				fmt.Println(token)
			}
		}
		`,
	}

	expectedFileContents := `
	package test

	func (p *Person) SayHello() {
		fmt.Println("hello")
	}
	func PrintToken() {
		tokens := []token.Token{
			token.VAR,
			token.CONST,
		}
	
		for _, token := range tokens {
			fmt.Println(token)
		}
	}
	`

	mockFileSet := token.NewFileSet()

	var files []*ast.File

	for _, mockFile := range mockFiles {
		file, err := parser.ParseFile(mockFileSet, "", mockFile, 0)
		assert.NoError(err)

		files = append(files, file)
	}

	concatFile, err := ConcatFiles(files, mockFileSet, nil)
	assert.NoError(err)

	var outPut bytes.Buffer
	if err := format.Node(&outPut, mockFileSet, concatFile); err != nil {
		assert.NoError(err)
	}

	expectedFile, err := parser.ParseFile(mockFileSet, "", expectedFileContents, 0)
	assert.NoError(err)

	var expectedOutput bytes.Buffer
	if err := format.Node(&expectedOutput, mockFileSet, expectedFile); err != nil {
		assert.NoError(err)
	}

	assert.Equal(expectedOutput.String(), outPut.String())
}

func TestDeleteFiles(t *testing.T) {
	assert := assert.New(t)
	mockContent := `package main`

	err := ioutil.WriteFile("test.go", []byte(mockContent), os.ModePerm)
	assert.NoError(err)

	mockFilePath := []string{"test.go"}

	DeleteFiles(mockFilePath)

	_, _, err = ParseASTFiles(mockFilePath)
	assert.Error(err)
}

func TestGetFilesToSort(t *testing.T) {
	assert := assert.New(t)

	mockFiles := []string{`
		package test

		const (
			name = "name"
		)
		`,
		`
		package test

		const (
			age = 30
		)
		`,
	}

	mockFileSet := token.NewFileSet()

	var files []*ast.File

	for _, mockFile := range mockFiles {
		file, err := parser.ParseFile(mockFileSet, "", mockFile, 0)
		assert.NoError(err)

		files = append(files, file)
	}

	options := NewOptions()

	sortedFiles, err := getFilesToSort(files, options, mockFileSet)
	assert.NoError(err)

	var output bytes.Buffer
	err = format.Node(&output, mockFileSet, sortedFiles[0])
	assert.NoError(err)

	assert.Equal("package goconcat\n\nconst (\n\tname = \"name\"\n\tage  = 30\n)\n", output.String())

}

func TestGetFilesToSort_ConcatPackages(t *testing.T) {
	assert := assert.New(t)

	mockFiles := []string{`
		package test

		const (
			name = "name"
		)
		`,
		`
		package test

		const (
			age = 30
		)
		`,
	}

	mockFileSet := token.NewFileSet()

	var files []*ast.File

	for _, mockFile := range mockFiles {
		file, err := parser.ParseFile(mockFileSet, "", mockFile, 0)
		assert.NoError(err)

		files = append(files, file)
	}

	options := NewOptions()
	options.SplitFilesByPackage = true

	sortedFiles, err := getFilesToSort(files, options, mockFileSet)
	assert.NoError(err)

	var output bytes.Buffer
	err = format.Node(&output, mockFileSet, sortedFiles[0])
	assert.NoError(err)

	assert.Equal("package test\n\nconst (\n\tname = \"name\"\n\tage  = 30\n)\n", output.String())

}

func TestGetFilePaths(t *testing.T) {
	assert := assert.New(t)
	options := NewOptions()

	options.SetOptions(
		".",
		[]Directory{},
		[]PrefixType{"file"},
		"",
		false,
		false,
	)

	paths, err := GetFilePaths(options)
	assert.NoError(err)

	assert.Equal([]string{"file.go", "file_test.go"}, paths)
}

func TestParseASTFiles(t *testing.T) {
	assert := assert.New(t)

	mockFilesContent := `package test
	var hello string
	`
	mockFileSet := token.NewFileSet()
	file, err := parser.ParseFile(mockFileSet, "", mockFilesContent, 0)
	assert.NoError(err)

	err = WriteASTFile(file, mockFileSet, "test.go")
	assert.NoError(err)

	mockFilePath := []string{"test.go"}

	parsedFile, parsedFileSet, err := ParseASTFiles(mockFilePath)
	assert.NoError(err)

	var output bytes.Buffer
	err = format.Node(&output, parsedFileSet, parsedFile[0])
	assert.NoError(err)

	assert.Equal("package test\n\nvar hello string\n", output.String())

	DeleteFiles(mockFilePath)
}

func TestCheckDirectoryIgnore(t *testing.T) {
	assert := assert.New(t)

	ignoredDirectories := []Directory{
		"test.go",
	}

	isValid := checkDirectoryIgnore("test.go", ignoredDirectories)

	assert.True(isValid)
}
