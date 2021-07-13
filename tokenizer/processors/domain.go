package processors

import (
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
)

type Walker func(root renders.Node, previous renders.Node, rootFun ast.Node)

type Expression func(exp ast.Expr) renders.Expr

type Processor interface {
	Process(node ast.Node, root renders.Node, previous renders.Node) error
}

type ProcessorCreator func(walk Walker, expression Expression) (Processor, error)

type Base struct {
	Walker     Walker
	Expression Expression
}

type ExpressionProcessorCreator func(expression Expression) (ExpressionProcessor, error)

type ExpressionProcessor interface {
	Process(exp ast.Expr) (renders.Expr, error)
}
