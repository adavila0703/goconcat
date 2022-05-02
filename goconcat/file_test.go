package goconcat

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"mockconcat/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	concatFile, err := ConcatFiles(files, mockFileSet)
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

	concatFile, err := ConcatFiles(files, mockFileSet)
	assert.NoError(err)

	var outPut bytes.Buffer
	if err := format.Node(&outPut, mockFileSet, concatFile); err != nil {
		assert.NoError(err)
	}

	expectedFile, err := parser.ParseFile(mockFileSet, "", expectedFileContents, 0)
	fmt.Println(err)
	assert.NoError(err)

	var expectedOutput bytes.Buffer
	if err := format.Node(&expectedOutput, mockFileSet, expectedFile); err != nil {
		assert.NoError(err)
	}

	assert.Equal(expectedOutput.String(), outPut.String())
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

	concatFile, err := ConcatFiles(files, mockFileSet)
	assert.NoError(err)

	var outPut bytes.Buffer
	if err := format.Node(&outPut, mockFileSet, concatFile); err != nil {
		assert.NoError(err)
	}

	expectedFile, err := parser.ParseFile(mockFileSet, "", expectedFileContents, 0)
	fmt.Println(err)
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

	concatFile, err := ConcatFiles(files, mockFileSet)
	assert.NoError(err)

	var outPut bytes.Buffer
	if err := format.Node(&outPut, mockFileSet, concatFile); err != nil {
		assert.NoError(err)
	}

	expectedFile, err := parser.ParseFile(mockFileSet, "", expectedFileContents, 0)
	fmt.Println(err)
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

	concatFile, err := ConcatFiles(files, mockFileSet)
	assert.NoError(err)

	var outPut bytes.Buffer
	if err := format.Node(&outPut, mockFileSet, concatFile); err != nil {
		assert.NoError(err)
	}

	expectedFile, err := parser.ParseFile(mockFileSet, "", expectedFileContents, 0)
	fmt.Println(err)
	assert.NoError(err)

	var expectedOutput bytes.Buffer
	if err := format.Node(&expectedOutput, mockFileSet, expectedFile); err != nil {
		assert.NoError(err)
	}

	assert.Equal(expectedOutput.String(), outPut.String())
}

func TestGetFilePaths(t *testing.T) {
	assert := assert.New(t)
	filePaths, err := GetFilePaths(
		".",
		[]utils.Directory{
			utils.DirectoryGit,
		},
		[]utils.FileType{
			utils.FileGo,
		},
		[]utils.PrefixType{
			utils.PrefixMock,
		},
	)
	assert.NoError(err)

	fmt.Println(filePaths)
}
