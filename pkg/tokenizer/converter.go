package tokenizer

import (
	"go/ast"
)

func FieldListToMap(params *ast.FieldList, t Transformer) map[string]string {
	res := map[string]string{}
	for _, param := range params.List {
		res[GetName(param.Names)] = t.Expressions(param.Type).String()
	}
	return res
}

func GetName(i []*ast.Ident) string {
	n := ""
	for _, name := range i {
		n += " " + name.String()
	}
	return n
}
