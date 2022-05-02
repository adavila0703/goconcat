package goconcat

import (
	"go/ast"
	"go/token"
	"io/fs"
	"mockconcat/utils"
	"path/filepath"
	"strings"
)

// concatenates all ast files
func ConcatFiles(files []*ast.File, fileSet *token.FileSet) (*ast.File, error) {
	if len(files) > 1 {
		targetFile, leftOverFiles := removeASTFileByIndex(files, 0)

		// concat targetfile
		concatenateTargetFile(targetFile)

		var importsToAdd []string

		for _, file := range leftOverFiles {
			for _, fileImports := range file.Imports {
				importsToAdd = append(importsToAdd, fileImports.Path.Value)
			}
		}

		ConcatImports(targetFile, fileSet, importsToAdd)

		// concat all types
		for _, file := range leftOverFiles {
			var tok token.Token

			tok = token.VAR
			spec, _ := GetSpecsAndIndices(file, tok)
			AddSpecToTargetFile(targetFile, spec, tok)

			tok = token.CONST
			spec, _ = GetSpecsAndIndices(file, tok)
			AddSpecToTargetFile(targetFile, spec, tok)

			tok = token.TYPE
			spec, _ = GetSpecsAndIndices(file, tok)
			AddSpecToTargetFile(targetFile, spec, tok)
		}

		funcs := GetFuncDeclFromFiles(leftOverFiles)

		targetFile.Decls = append(targetFile.Decls, funcs...)

		return targetFile, nil
	}
	return nil, ErrReadingDirectories
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
	specs, indices := GetSpecsAndIndices(file, tok)

	specs = utils.RemoveFromSlice(specs, 0)

	RemoveDecl(file, indices)

	ConcatSpecs(file, specs, tok)
}

func GetFilePaths(
	path string,
	ignoredDirectories []utils.Directory,
	fileTypes []utils.FileType,
	prefix []utils.PrefixType,
) ([]string, error) {
	var filePaths []string

	fileTypeMap := make(map[utils.FileType]utils.FileType)

	for _, fileType := range fileTypes {
		if _, ok := fileTypeMap[fileType]; ok {
			continue
		}
		fileTypeMap[fileType] = fileType
	}

	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if checkDirectoryIgnore(path, ignoredDirectories) {
			return nil
		}

		if !info.IsDir() {
			suffix := getSuffixFileType(info.Name())
			if _, ok := fileTypeMap[suffix]; !ok {
				return nil
			}

			hasPrefix := containsPrefix(info, prefix)

			if !hasPrefix {
				return nil
			}

			filePaths = append(filePaths, path)
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

func checkDirectoryIgnore(path string, ignoredDirectories []utils.Directory) bool {
	for _, ignoreDirectory := range ignoredDirectories {
		directory := utils.AnyToString(ignoreDirectory)
		if strings.Contains(path, directory) {
			return true
		}
	}
	return false
}

type Specs interface {
	[]ast.Spec
}

func AddSpecToTargetFile[T Specs](targetFile *ast.File, spec T, tok token.Token) {
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

func getSuffixFileType(fileName string) utils.FileType {
	return utils.FileType(filepath.Ext(fileName))
}

func containsPrefix(info fs.FileInfo, prefix []utils.PrefixType) bool {
	for _, p := range prefix {
		sPrefix := utils.AnyToString(p)
		if strings.HasPrefix(info.Name(), sPrefix) {
			return true
		}
	}
	return false
}
