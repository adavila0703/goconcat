package goconcat

import (
	"go/ast"
	"go/token"
)

func ConcatImports(targetFile *ast.File, fileSet *token.FileSet, importStrings []string) {
	existingImports := make(map[string]string)

	for _, v := range targetFile.Imports {
		existingImports[v.Path.Value] = v.Path.Value
	}

	for _, importString := range importStrings {
		// skip import if it already exists
		if _, ok := existingImports[importString]; ok {
			continue
		}
		addImportToTargetFile(targetFile, importString)
	}
	ast.SortImports(fileSet, targetFile)
}

func addImportToTargetFile(targetFile *ast.File, target string) {
	for _, decl := range targetFile.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)

			if genDecl.Tok == token.IMPORT {
				spec := &ast.ImportSpec{
					Path: &ast.BasicLit{
						Value: target,
					},
				}
				genDecl.Specs = append(genDecl.Specs, spec)
			}
		}
	}
}
