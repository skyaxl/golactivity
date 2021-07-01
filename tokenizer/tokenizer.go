package tokenizer

import (
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/drawer"
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
			if fi.Init != nil {
				fiNo.Init = &drawer.Assignation{
					BaseNode: drawer.BaseNode{
						Par: fiNo,
						Dep: root.Depth() + 1,
					},
					Left:  make([]drawer.Expr, 0),
					Right: make([]drawer.Expr, 0),
				}
				if ass, ok := fi.Init.(*ast.AssignStmt); ok {
					for _, ex := range ass.Lhs {
						fiNo.Init.Left = append(fiNo.Init.Left, t.Expressions(ex))
					}
					for _, ex := range ass.Rhs {
						fiNo.Init.Right = append(fiNo.Init.Right, t.Expressions(ex))
					}
				}
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
	case *ast.AssignStmt:
		{
			as := rootFun.(*ast.AssignStmt)
			ass := &drawer.Assignation{
				BaseNode: drawer.BaseNode{
					Par: root,
					Dep: root.Depth() + 1,
				},
				Left:  make([]drawer.Expr, 0),
				Right: make([]drawer.Expr, 0),
			}
			for _, ex := range as.Lhs {
				ass.Left = append(ass.Left, t.Expressions(ex))
			}
			for _, ex := range as.Rhs {
				ass.Right = append(ass.Right, t.Expressions(ex))
			}

			if previous != nil {
				previous.SetNext(ass)
			}
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
				Exp: t.Expressions(call),
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
	case *ast.ForStmt:
		{
			fors := rootFun.(*ast.ForStmt)
			forr := &drawer.For{
				BaseNode: drawer.BaseNode{
					Par:  root,
					Prev: previous,
					Dep:  root.Depth() + 1,
				},
				Conditions: t.Expressions(fors.Cond),
			}
			if fors.Init != nil {
				forr.Init = &drawer.Assignation{
					BaseNode: drawer.BaseNode{
						Par: forr,
						Dep: root.Depth() + 1,
					},
					Left:  make([]drawer.Expr, 0),
					Right: make([]drawer.Expr, 0),
				}
				if ass, ok := fors.Init.(*ast.AssignStmt); ok {
					for _, ex := range ass.Lhs {
						forr.Init.Left = append(forr.Init.Left, t.Expressions(ex))
					}
					for _, ex := range ass.Rhs {
						forr.Init.Right = append(forr.Init.Right, t.Expressions(ex))
					}
				}
			}

			if fors.Post != nil {

			}

			forr.Body = &drawer.Root{}
			forr.Body.Par = forr

			t.Walk(forr, forr.Body, fors.Body)
			if previous != nil {
				previous.SetNext(forr)
			}
			return
		}
	case *ast.RangeStmt:
		{
			fors := rootFun.(*ast.RangeStmt)
			forr := &drawer.Range{
				BaseNode: drawer.BaseNode{
					Par:  root,
					Prev: previous,
					Dep:  root.Depth() + 1,
				},
				ID:    t.Expressions(fors.X),
				Key:   t.Expressions(fors.Key),
				Value: t.Expressions(fors.Value),
				Body:  &drawer.Root{},
			}

			forr.Body.Par = forr
			t.Walk(forr, forr.Body, fors.Body)
			if previous != nil {
				previous.SetNext(forr)
			}
			return
		}
	case *ast.SwitchStmt:
		{
			sw := rootFun.(*ast.SwitchStmt)
			swi := &drawer.Switch{
				BaseNode: drawer.BaseNode{
					Par:  root,
					Prev: previous,
					Dep:  root.Depth() + 1,
				},
				Cases: make([]*drawer.Case, 0),
				Tag:   t.Expressions(sw.Tag),
			}
			if sw.Init != nil {
				swi.Init = &drawer.Assignation{
					BaseNode: drawer.BaseNode{
						Par: swi,
						Dep: root.Depth(),
					},
					Left:  make([]drawer.Expr, 0),
					Right: make([]drawer.Expr, 0),
				}
				if ass, ok := sw.Init.(*ast.AssignStmt); ok {
					for _, ex := range ass.Lhs {
						swi.Init.Left = append(swi.Init.Left, t.Expressions(ex))
					}
					for _, ex := range ass.Rhs {
						swi.Init.Right = append(swi.Init.Right, t.Expressions(ex))
					}
				}
			}

			if sw.Body != nil {
				for _, c := range sw.Body.List {
					cas := c.(*ast.CaseClause)
					swic := &drawer.Case{
						BaseNode: drawer.BaseNode{
							Par: swi,
							Dep: swi.Depth(),
						},
						Value: make(drawer.Expressions, 0),
						Body:  &drawer.Root{},
					}
					swic.Body.BaseNode = drawer.BaseNode{
						Par: swic,
						Dep: swic.Depth(),
					}
					if cas.List != nil {
						for _, exp := range cas.List {
							swic.Value = append(swic.Value, t.Expressions(exp))
						}
					}

					t.Walk(swic.Body, nil, &ast.BlockStmt{
						List: cas.Body,
					})
					swi.Cases = append(swi.Cases, swic)
				}

			}

			if previous != nil {
				previous.SetNext(swi)
			}
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
	case *ast.ChanType:
		{
			u := exp.(*ast.ChanType)
			v := drawer.Chan{}
			v.Value = t.Expressions(u.Value)
			return v
		}
	case *ast.FuncLit:
		{
			u := exp.(*ast.FuncLit)
			v := drawer.FunLiteral{
				Args:      make(drawer.Expressions, 0),
				Responses: make(drawer.Expressions, 0),
			}
			for _, a := range u.Type.Params.List {
				v.Args = append(v.Args, drawer.Field{
					Name: GetName(a.Names),
					Kind: t.Expressions(a.Type),
				})
			}
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
	case *ast.CallExpr:
		{
			u := exp.(*ast.CallExpr)
			v := drawer.Call{
				Func:      t.Expressions(u.Fun).(drawer.Identifier),
				Arguments: make(drawer.Expressions, 0),
			}
			for _, arg := range u.Args {
				v.Arguments = append(v.Arguments, t.Expressions(arg))
			}

			return v
		}
	case *ast.CompositeLit:
		{
			u := exp.(*ast.CompositeLit)
			v := drawer.Literal{
				Kind:     t.Expressions(u.Type),
				Elements: make(drawer.Expressions, 0),
			}

			for _, arg := range u.Elts {
				v.Elements = append(v.Elements, t.Expressions(arg))
			}
			return v
		}
	case *ast.ArrayType:
		{
			u := exp.(*ast.ArrayType)
			v := drawer.ArrayType{
				Type: t.Expressions(u.Elt),
			}
			if u.Len != nil {
				v.Len = t.Expressions(u.Len)
			}
			return v
		}
	case *ast.IndexExpr:
		{
			u := exp.(*ast.IndexExpr)
			v := drawer.Index{
				Ident: t.Expressions(u.X),
				Index: t.Expressions(u.Index),
			}
			return v
		}
	}

	return nil
}
