package finder

import (
	"go/ast"
	"go/token"
	"strings"
)

var (
	// DrawAspect defines the comment that should be in top of method to be draw
	DrawAspect = "@draw"
)

//ExistsComment verifies if exist aspect comment
func ExistsComment(comments []*ast.Comment, aspect string) bool {
	for _, comment := range comments {
		if strings.Contains(comment.Text, string(aspect)) {
			return true
		}
	}
	return false
}

// FindAllowedFuncs finds the funcs that have the selected aspect setted
func FindAllowedFuncs(pkgs map[string]*ast.Package, fset *token.FileSet) []*ast.FuncDecl {
	finals := []*ast.FuncDecl{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			funs := getFileFuncs(file)
			finals = append(finals, funs...)
		}
	}
	return finals
}

func getFileFuncs(pkg *ast.File) (res []*ast.FuncDecl) {
	res = make([]*ast.FuncDecl, 0)
	for _, decs := range pkg.Decls {
		switch dcs := decs.(type) {
		case *ast.FuncDecl:
			{
				f := (*ast.FuncDecl)(dcs)
				if f.Doc == nil || !ExistsComment(f.Doc.List, DrawAspect) {
					break
				}
				res = append(res, f)
			}
		}
	}
	return res
}
