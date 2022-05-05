package goconcat

import (
	"go/ast"
)

// removes decl from ast file
func RemoveDecl(file *ast.File, indices []int) {
	// get the first index of decl as the base decl
	indices = RemoveFromSlice(indices, 0)
	file.Decls = ReturnAllButIndices(file.Decls, indices)
}

func GetFuncDeclFromFiles(files []*ast.File) []ast.Decl {
	var funcs []ast.Decl

	for _, file := range files {
		for _, decl := range file.Decls {
			switch decl.(type) {
			case *ast.FuncDecl:
				funcDecl := decl.(*ast.FuncDecl)
				funcs = append(funcs, funcDecl)
			}
		}
	}

	return funcs
}
