package goconcat

import (
	"go/ast"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

// concatenates all ast files
func ConcatFiles(files []*ast.File, fileSet *token.FileSet) (*ast.File, error) {
	if len(files) > 1 {
		targetFile, leftOverFiles := removeASTFileByIndex(files, 0)

		var importsToAdd []string

		for _, file := range leftOverFiles {
			for _, fileImports := range file.Imports {
				importsToAdd = append(importsToAdd, fileImports.Path.Value)
			}
		}

		ConcatImports(targetFile, fileSet, importsToAdd)

		for _, file := range leftOverFiles {
			spec := GetVarGenDecl(file)
			AddSpecToTargetFile(targetFile, spec, token.VAR)
		}

		return targetFile, nil
	}
	return nil, ErrReadingDirectories
}

func GetFilePaths(path string, ignoredDirectories []string, fileType string, prefix string) ([]string, error) {
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

func checkDirectoryIgnore(path string, ignoredDirectories []string) bool {
	for _, ignoreDirectory := range ignoredDirectories {
		if strings.Contains(path, ignoreDirectory) {
			return true
		}
	}
	return false
}

type test interface {
	[]ast.Spec
}

func AddSpecToTargetFile[T test](targetFile *ast.File, spec T, tok token.Token) {
	switch tok {
	case token.VAR:
		for _, decl := range targetFile.Decls {
			switch decl.(type) {
			case *ast.GenDecl:
				genDecl := decl.(*ast.GenDecl)
				if genDecl.Tok == token.VAR {
					genDecl.Specs = append(genDecl.Specs, spec...)
				}
			}
		}
	case token.IMPORT:
	}
}
