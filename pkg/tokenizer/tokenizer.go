package tokenizer

import (
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/pkg/drawer"
)

type Transformer struct {
	funcs *ast.FuncDecl
}

//NewTransformer new transformer
func NewTransformer(funcs *ast.FuncDecl) *Transformer {
	return &Transformer{funcs: funcs}
}

//Transform
func (t Transformer) Transform() (doc *drawer.Document) {
	fun := t.funcs
	doc = &drawer.Document{
		Comment: fun.Doc.Text(),
		Name:    fun.Name.Name,
	}
	doc.Root = &drawer.Root{}
	doc.Root.Params = FieldListToMap(fun.Type.Params, t)
	if fun.Type.Results != nil {
		doc.Root.Responses = FieldListToMap(fun.Type.Results, t)
	}
	t.Walk(doc.Root, nil, fun.Body)

	return doc
}

func (t Transformer) Walk(root drawer.Node, previous drawer.Node, rootFun ast.Node) {
	if previous == nil {
		previous = root
	}
	switch rootFun.(type) {
	case *ast.FuncDecl:
		{
			panic("Func declaration is not implemented")
		}
	case *ast.BlockStmt:
		{
			b := rootFun.(*ast.BlockStmt)

			for _, n := range b.List {
				t.Walk(root, previous, n)
				if previous != nil && previous.Next() != nil {
					previous = previous.Next()
				}
			}
			return
		}
	case *ast.IfStmt:
		{
			fi := rootFun.(*ast.IfStmt)
			fiNo := &drawer.If{
				BaseNode: drawer.BaseNode{
					Par:  root,
					Prev: previous,
					Dep:  root.Depth() + 1,
				},
				Conditions: t.Expressions(fi.Cond),
			}

			fiNo.Body = &drawer.Root{}
			fiNo.Body.Par = fiNo

			t.Walk(fiNo, fiNo.Body, fi.Body)
			if previous != nil {
				previous.SetNext(fiNo)
			}

			return
		}
	case *ast.ExprStmt:
		{
			call := rootFun.(*ast.ExprStmt)
			t.Walk(root, previous, call.X)
		}
	case *ast.CallExpr:
		{
			call := rootFun.(*ast.CallExpr)
			act := &drawer.Activity{
				BaseNode: drawer.BaseNode{
					Par:  root,
					Prev: previous,
					Dep:  root.Depth() + 1,
				},
				Name: t.Expressions(call.Fun).String(),
			}
			if previous != nil {
				previous.SetNext(act)
			}
			return
		}
	case *ast.ReturnStmt:
		{
			ret := rootFun.(*ast.ReturnStmt)
			act := &drawer.Return{
				BaseNode: drawer.BaseNode{
					Par:  root,
					Prev: previous,
					Dep:  root.Depth() + 1,
				},
				Values: []drawer.Expr{},
			}

			for _, exp := range ret.Results {
				act.Values = append(act.Values, t.Expressions(exp))
			}

			if previous != nil {
				previous.SetNext(act)
			}
		}
	case *ast.SwitchStmt:
		{
			panic("Not inplemented")
		}

	}
}

//Expressions transform expressions
// *ast.Binary, *ast.Unary, *ast.Ident, *ast.BasicLit
func (t Transformer) Expressions(exp ast.Expr) drawer.Expr {
	switch exp.(type) {
	case *ast.UnaryExpr:
		{
			u := exp.(*ast.UnaryExpr)
			du := drawer.Unary{}
			du.Oper = u.Op.String()
			du.Left = t.Expressions(u.X)
			return du
		}
	case *ast.Ident:
		{
			u := exp.(*ast.Ident)
			du := drawer.Identifier{}
			du.ID = u.String()
			return du
		}
	case *ast.BinaryExpr:
		{
			u := exp.(*ast.BinaryExpr)
			du := drawer.Binary{}
			du.Oper = u.Op.String()
			du.Left = t.Expressions(u.X)
			du.Right = t.Expressions(u.Y)
			return du
		}
	case *ast.ParenExpr:
		{
			u := exp.(*ast.ParenExpr)
			du := drawer.Parent{}
			du.Expr = t.Expressions(u.X)
			return du
		}
	case *ast.BasicLit:
		{
			u := exp.(*ast.BasicLit)
			v := drawer.Value{}
			v.Value = u.Value
			v.Kind = u.Kind.String()
			return v
		}
	case *ast.SelectorExpr:
		{
			u := exp.(*ast.SelectorExpr)
			v := drawer.Identifier{}
			x := t.Expressions(u.X)
			sel := t.Expressions(u.Sel)
			v.ID = fmt.Sprintf("%s.%s", x.String(), sel.String())
			return v
		}
	}

	return nil
}
