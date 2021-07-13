package expressions

import "go/ast"

func GetName(i []*ast.Ident) string {
	n := ""
	for _, name := range i {
		n += " " + name.String()
	}
	return n
}
