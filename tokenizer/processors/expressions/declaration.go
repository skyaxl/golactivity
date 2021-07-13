package expressions

import (
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type Chan struct {
	processors.Base
}

func NewChan(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &Chan{processors.Base{Expression: exp}}, nil
}

func (t *Chan) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.ChanType)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.BasicLit", exp)
	}

	v := renders.Chan{}
	v.Value = t.Expression(u.Value)
	return v, nil
}

type Selector struct {
	processors.Base
}

func NewSelector(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &Selector{processors.Base{Expression: exp}}, nil
}

func (t *Selector) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.SelectorExpr)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.SelectorExpr", exp)
	}
	v := renders.Identifier{}
	x := t.Expression(u.X)
	sel := t.Expression(u.Sel)
	v.ID = fmt.Sprintf("%s.%s", x.String(), sel.String())
	return v, nil
}

type ArrayType struct {
	processors.Base
}

func NewArrayType(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &ArrayType{processors.Base{Expression: exp}}, nil
}

func (t *ArrayType) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.ArrayType)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.SelectorExpr", exp)
	}
	v := renders.ArrayType{
		Type: t.Expression(u.Elt),
	}
	if u.Len != nil {
		v.Len = t.Expression(u.Len)
	}
	return v, nil
}

type Index struct {
	processors.Base
}

func NewIndex(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &Index{processors.Base{Expression: exp}}, nil
}

func (t *Index) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.IndexExpr)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.SelectorExpr", exp)
	}
	v := renders.Index{
		Ident: t.Expression(u.X),
		Index: t.Expression(u.Index),
	}
	return v, nil
}
