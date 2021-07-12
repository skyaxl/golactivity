package blocks

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type Expr struct {
	processors.Base
}

// NewExprProcessor creates new expression processor
func NewExprProcessor(walk processors.Walker, expression processors.Expression) (processors.Processor, error) {
	if walk == nil || expression == nil {
		return nil, errors.New("invalid parameters")
	}
	return &Expr{Base: processors.Base{Walker: walk, Expression: expression}}, nil
}

func (t *Expr) Process(node ast.Node, root renders.Node, previous renders.Node) error {
	call, ok := node.(*ast.ExprStmt)
	if !ok {
		return fmt.Errorf("The node %v is not a *ast.ReturnStmt", node)
	}
	t.Walker(root, previous, call.X)
	return nil
}

type Assignation struct {
	processors.Base
}

// NewAssignationProcessor creates new assignation processor
func NewAssignationProcessor(walk processors.Walker, expression processors.Expression) (processors.Processor, error) {
	if walk == nil || expression == nil {
		return nil, errors.New("invalid parameters")
	}
	return &Assignation{Base: processors.Base{Walker: walk, Expression: expression}}, nil
}

func (t *Assignation) Process(node ast.Node, root renders.Node, previous renders.Node) error {
	as, ok := node.(*ast.AssignStmt)
	if !ok {
		return fmt.Errorf("The node %v is not a *ast.ReturnStmt", node)
	}
	ass := &renders.Assignation{
		BaseNode: renders.BaseNode{
			Par: root,
			Dep: root.Depth() + 1,
		},
		Left:  make([]renders.Expr, 0),
		Right: make([]renders.Expr, 0),
	}
	for _, ex := range as.Lhs {
		ass.Left = append(ass.Left, t.Expression(ex))
	}
	for _, ex := range as.Rhs {
		ass.Right = append(ass.Right, t.Expression(ex))
	}

	if previous != nil {
		previous.SetNext(ass)
	}
	return nil
}
