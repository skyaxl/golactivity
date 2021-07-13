package expressions

import (
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type Unary struct {
	processors.Base
}

func NewUnary(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &Unary{processors.Base{Expression: exp}}, nil
}

func (t *Unary) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.UnaryExpr)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.UnaryExpr", exp)
	}
	du := renders.Unary{}
	du.Oper = u.Op.String()
	du.Left = t.Expression(u.X)
	return du, nil
}

type Binary struct {
	processors.Base
}

func NewBinary(exp processors.Expression) (processors.ExpressionProcessor, error) {
	if exp == nil {
		return nil, fmt.Errorf("Invalid parameters")
	}
	return &Binary{processors.Base{Expression: exp}}, nil
}

func (t *Binary) Process(exp ast.Expr) (renders.Expr, error) {
	u, ok := exp.(*ast.BinaryExpr)
	if !ok {
		return nil, fmt.Errorf("The expression %v is not a *ast.BinaryExpr", exp)
	}
	du := renders.Binary{}
	du.Oper = u.Op.String()
	du.Left = t.Expression(u.X)
	du.Right = t.Expression(u.Y)
	return du, nil
}
