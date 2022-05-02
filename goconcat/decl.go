package goconcat

import (
	"go/ast"
	"mockconcat/utils"
)

// removes decl from ast file
func RemoveDecl(file *ast.File, indices []int) {
	// get the first index of decl as the base decl
	indices = utils.RemoveFromSlice(indices, 0)
	file.Decls = utils.ReturnAllButIndices(file.Decls, indices)
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
