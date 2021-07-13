package expressions

import (
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type Ident struct {
	processors.Base
}

func NewIdent(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &Ident{processors.Base{Expression: exp}}, nil
}

func (t *Ident) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.Ident", exp)
	}
	du := renders.Identifier{}
	du.ID = u.String()
	return du, nil
}

type Paren struct {
	processors.Base
}

func NewParen(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &Paren{processors.Base{Expression: exp}}, nil
}

func (t *Paren) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.ParenExpr)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.BinaryExpr", exp)
	}
	du := renders.Parent{}
	du.Expr = t.Expression(u.X)
	return du, nil
}

type BasicLit struct {
	processors.Base
}

func NewBasicLit(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &BasicLit{processors.Base{Expression: exp}}, nil
}

func (t *BasicLit) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.BasicLit)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.BasicLit", exp)
	}
	v := renders.Value{}
	v.Value = u.Value
	v.Kind = u.Kind.String()
	return v, nil
}

type CompositeLit struct {
	processors.Base
}

func NewCompositeLit(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &CompositeLit{processors.Base{Expression: exp}}, nil
}

func (t *CompositeLit) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.SelectorExpr", exp)
	}
	v := renders.Literal{
		Kind:     t.Expression(u.Type),
		Elements: make(renders.Expressions, 0),
	}

	for _, arg := range u.Elts {
		v.Elements = append(v.Elements, t.Expression(arg))
	}
	return v, nil
}
