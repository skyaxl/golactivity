package tokenizer

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

func ReadTokens(pkgs map[string]*ast.Package, fset *token.FileSet) []*ast.FuncDecl {
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
	res = make([]*ast.FuncDecl, 0, 0)
	for _, decs := range pkg.Decls {
		switch decs.(type) {
		case *ast.FuncDecl:
			{
				f := decs.(*ast.FuncDecl)
				if f.Doc == nil || !ExistsComment(f.Doc.List, DrawAspect) {
					break
				}
				res = append(res, f)
			}
		}
	}
	return res
}

// Properties to inspect
//Body,List,Condition
func inspectBody(block []ast.Stmt, node string, depth int) {
	for _, dec := range block {
		fmt.Printf("%sNode %s, type: %s, %+1v\n", strings.Repeat("\t", depth), node, reflect.ValueOf(dec).Elem().Type().Name(), dec)
		processBody(dec, depth+1)
		processList(dec, depth+1)
		processExpression(dec, depth+1)
		processCondition(dec, depth+1)
		processInnerExpression(dec, depth+1)
	}
}

func inspectBodyInterface(block []interface{}, node string, depth int) {
	for _, dec := range block {
		fmt.Printf("Node %s, type: %s, %+1v\n", node, reflect.ValueOf(dec).Elem().Type().Name(), dec)
		processBody(dec, depth+1)
		processList(dec, depth+1)
		processExpression(dec, depth+1)
		processCondition(dec, depth+1)
		processInnerExpression(dec, depth+1)
	}
}

func processExpression(dec interface{}, depth int) {
	var (
		e  *ast.ExprStmt
		ok bool
	)
	if e, ok = dec.(*ast.ExprStmt); !ok {
		return
	}

	fmt.Printf("%sNode %s, type: %s, %+1v\n", strings.Repeat("\t", depth), "expression", reflect.ValueOf(dec).Elem().Type().Name(), e.X)
	processInnerExpression(e.X, depth+1)
}

func processInnerExpression(exp interface{}, depth int) {

	fmt.Printf("%sNode %s, type: %s, %+1v\n", strings.Repeat("\t", depth), "inner_expression", reflect.ValueOf(exp).Elem().Type().Name(), exp)
	v := reflect.ValueOf(exp).Elem()
	if !v.CanInterface() || !v.CanAddr() || v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		t := field.Type()
		implements := (t.Kind() == reflect.Interface && t.String() == "ast.Expr")
		if !implements || field.IsNil() || !field.CanInterface() || field.Elem().IsNil() {
			continue
		}
		processInnerExpression(field.Interface().(ast.Expr), depth+1)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		t := field.Type()
		implements := (t.Kind() == reflect.Slice && t.String() == "[]ast.Expr")
		if !implements || field.IsNil() || !field.CanInterface() {
			continue
		}
		express := field.Interface().([]ast.Expr)
		for _, e := range express {
			processInnerExpression(e, depth+1)
		}
	}
}

func processCondition(dec interface{}, depth int) {
	v := reflect.ValueOf(dec).Elem()
	list := v.FieldByName("Cond")

	if !list.CanAddr() || !list.CanInterface() || list.IsNil() {
		return
	}
	cond, _ := list.Interface().(ast.Expr)
	if cond == nil {
		return
	}
	processInnerExpression(cond, depth+1)
}

func processList(dec interface{}, depth int) {
	v := reflect.ValueOf(dec).Elem()
	list := v.FieldByName("List")

	if !list.CanAddr() || !list.CanInterface() || list.IsNil() || list.Kind() != reflect.Slice {
		return
	}
	slice, _ := list.Interface().([]ast.Stmt)
	if slice == nil {
		return
	}
	inspectBody(slice, "list", depth+1)
}

func processBody(dec interface{}, depth int) {
	v := reflect.ValueOf(dec).Elem()
	body := v.FieldByName("Body")
	if body.CanAddr() && body.CanInterface() && !body.IsNil() {
		b, _ := body.Interface().(*ast.BlockStmt)
		if b != nil && b.List != nil {
			inspectBody(b.List, "body", depth+1)
			return
		}
		if body.Kind() == reflect.Slice {
			a, _ := body.Interface().([]ast.Stmt)
			inspectBody(a, "body", depth+1)
			return
		}
		elem := body.Elem()
		if elem.CanAddr() && elem.CanInterface() && !elem.IsNil() {
			b, _ := elem.Interface().(*ast.BlockStmt)
			if b != nil && b.List != nil {
				inspectBody(b.List, "body.list", depth+1)
			}
		}
	}
}
