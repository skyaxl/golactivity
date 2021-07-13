package expressions

import (
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type FuncLit struct {
	processors.Base
}

func NewFuncLit(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &FuncLit{processors.Base{Expression: exp}}, nil
}

func (t *FuncLit) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.FuncLit)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.BasicLit", exp)
	}

	v := renders.FunLiteral{
		Args:      make(renders.Expressions, 0),
		Responses: make(renders.Expressions, 0),
	}
	for _, a := range u.Type.Params.List {
		v.Args = append(v.Args, renders.Field{
			Name: GetName(a.Names),
			Kind: t.Expression(a.Type),
		})
	}

	for _, a := range u.Type.Results.List {
		v.Responses = append(v.Responses, renders.Field{
			Name: GetName(a.Names),
			Kind: t.Expression(a.Type),
		})
	}
	return v, nil
}

type Call struct {
	processors.Base
}

func NewCall(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &Call{processors.Base{Expression: exp}}, nil
}

func (t *Call) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.CallExpr)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.SelectorExpr", exp)
	}
	v := renders.Call{
		Func:      t.Expression(u.Fun).(renders.Identifier),
		Arguments: make(renders.Expressions, 0),
	}
	for _, arg := range u.Args {
		v.Arguments = append(v.Arguments, t.Expression(arg))
	}

	return v, nil
}
