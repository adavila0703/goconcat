package concat

import (
	"go/ast"
	"go/token"
)

func GetSpecsAndIndices(file *ast.File, tok token.Token) ([]ast.Spec, []int) {
	var specs []ast.Spec
	var declIndex []int

	for index, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl.Tok == tok {
				declIndex = append(declIndex, index)
				specs = append(specs, genDecl.Specs...)
			}
		}
	}

	return specs, declIndex
}

func ConcatSpecs(file *ast.File, specs []ast.Spec, tok token.Token) {
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			if genDecl.Tok == tok {
				genDecl.Specs = append(genDecl.Specs, specs...)
			}
		}
	}
}
