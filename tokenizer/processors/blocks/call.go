package blocks

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type Call struct {
	processors.Base
}

// NewAssignationProcessor creates new assignation processor
func NewCallProcessor(walk processors.Walker, expression processors.Expression) (processors.Processor, error) {
	if walk == nil || expression == nil {
		return nil, errors.New("invalid parameters")
	}
	return &Call{Base: processors.Base{Walker: walk, Expression: expression}}, nil
}

func (t *Call) Process(node ast.Node, root renders.Node, previous renders.Node) error {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return fmt.Errorf("The node %v is not a *ast.ReturnStmt", node)
	}

	act := &renders.Activity{
		BaseNode: renders.BaseNode{
			Par:  root,
			Prev: previous,
			Dep:  root.Depth() + 1,
		},
		Exp: t.Expression(call),
	}
	if previous != nil {
		previous.SetNext(act)
	}
	return nil
}
