package goconcat

import (
	"go/ast"
	"go/token"
)

func GetVarGenDecl(file *ast.File) []ast.Spec {
	var spec []ast.Spec

	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl.Tok == token.VAR {
				spec = genDecl.Specs
			}
		}
	}

	return spec
}
