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
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// concatenates all ast files
func ConcatFiles(files []*ast.File, fileSet *token.FileSet) (*ast.File, error) {
	targetFile, leftOverFiles := removeASTFileByIndex(files, 0)
	if len(leftOverFiles) < 1 {
		return targetFile, nil
	}

	// concat targetfile
	concatenateTargetFile(targetFile)

	var importsToAdd []string

	for _, file := range leftOverFiles {
		for _, fileImports := range file.Imports {
			importsToAdd = append(importsToAdd, fileImports.Path.Value)
		}
	}

	concatImports(targetFile, fileSet, importsToAdd)

	// concat all types
	for _, file := range leftOverFiles {
		var tok token.Token

		tok = token.VAR
		spec, _ := getSpecsAndIndices(file, tok)
		addSpecToTargetFile(targetFile, spec, tok)

		tok = token.CONST
		spec, _ = getSpecsAndIndices(file, tok)
		addSpecToTargetFile(targetFile, spec, tok)

		tok = token.TYPE
		spec, _ = getSpecsAndIndices(file, tok)
		addSpecToTargetFile(targetFile, spec, tok)
	}

	funcs := getFuncDeclFromFiles(leftOverFiles)

	targetFile.Decls = append(targetFile.Decls, funcs...)

	return targetFile, nil
}

func GetFilePaths(options *Options) ([]string, error) {
	if err := validateOptions(options); err != nil {
		return nil, errors.WithStack(err)
	}

	var filePaths []string

	fileTypeMap := make(map[FileType]FileType)

	for _, fileType := range options.FileType {
		if _, ok := fileTypeMap[fileType]; ok {
			continue
		}
		fileTypeMap[fileType] = fileType
	}

	filepath.Walk(options.RootPath, func(path string, info fs.FileInfo, err error) error {
		if checkDirectoryIgnore(path, options.IgnoredDirectories) {
			return nil
		}

		if !info.IsDir() {
			suffix := getSuffixFileType(info.Name())
			if _, ok := fileTypeMap[suffix]; !ok {
				return nil
			}

			hasPrefix := containsPrefix(info, options.FilePrefix)

			if !hasPrefix {
				return nil
			}

			filePaths = append(filePaths, path)
		}
		return nil
	})

	if len(filePaths) < 1 {
		return nil, errors.WithStack(errNoFilesDetected)
	}

	return filePaths, nil
}

func DeleteFiles(filePaths []string) error {
	for _, file := range filePaths {
		if err := os.Remove(file); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// can return single ast file or list of files depending if you want to sort by package
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
		file, err := ConcatFiles(files, fileSet)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		filesToSort = append(filesToSort, file)
	}

	return filesToSort, nil
}

// parse a list of file paths to get ast files
func ParseASTFiles(filePaths []string) ([]*ast.File, *token.FileSet, error) {
	var filesToConcat []*ast.File
	fileSet := token.NewFileSet()

	if len(filePaths) < 1 {
		return nil, nil, errors.WithStack(errNoFilePath)
	}

	for _, path := range filePaths {
		fileContents, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}

		astFiles, err := parser.ParseFile(fileSet, "", fileContents, 0)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}

		filesToConcat = append(filesToConcat, astFiles)
	}

	return filesToConcat, fileSet, nil
}

func WriteASTFile(file *ast.File, fileSet *token.FileSet, filePath string) error {
	var buf bytes.Buffer
	if err := format.Node(&buf, fileSet, file); err != nil {
		return errors.WithStack(err)
	}
	if err := ioutil.WriteFile(filePath, buf.Bytes(), os.ModePerm); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func concatenateTargetFile(file *ast.File) {
	tokens := []token.Token{
		token.VAR,
		token.CONST,
		token.TYPE,
	}

	for _, token := range tokens {
		concatType(file, token)
	}
}

func concatType(file *ast.File, tok token.Token) {
	specs, indices := getSpecsAndIndices(file, tok)

	specs = removeFromSlice(specs, 0)

	removeDecl(file, indices)

	concatSpecs(file, specs, tok)
}

func removeASTFileByIndex(files []*ast.File, fileIndex int) (targetFile *ast.File, leftOverFiles []*ast.File) {
	for index, file := range files {
		if index == fileIndex {
			targetFile = file
			continue
		}
		leftOverFiles = append(leftOverFiles, file)
	}
	return
}

func checkDirectoryIgnore(path string, ignoredDirectories []Directory) bool {
	for _, ignoreDirectory := range ignoredDirectories {
		directory := anyToString(ignoreDirectory)
		if strings.Contains(path, directory) {
			return true
		}
	}
	return false
}

type specs interface {
	[]ast.Spec
}

func addSpecToTargetFile[T specs](targetFile *ast.File, spec T, tok token.Token) {
	for _, decl := range targetFile.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl.Tok == tok {
				genDecl.Specs = append(genDecl.Specs, spec...)
			}
		}
	}
}

func getSuffixFileType(fileName string) FileType {
	return FileType(filepath.Ext(fileName))
}

func containsPrefix(info fs.FileInfo, prefix []PrefixType) bool {
	if prefix == nil {
		return true
	}

	for _, p := range prefix {
		sPrefix := anyToString(p)
		if strings.HasPrefix(info.Name(), sPrefix) {
			return true
		}
	}
	return false
}

func destinationDirIsValid(rootPath string, destination string) bool {
	des := anyToString(destination)
	dirIsValid := false

	filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == des {
				dirIsValid = true
			}
		}
		return nil
	})

	return dirIsValid
}
