package control

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type Return struct {
	processors.Base
}

// NewReturnProcessor creates new return processor to generate renders.Return
func NewReturnProcessor(walk processors.Walker, expression processors.Expression) (processors.Processor, error) {
	if walk == nil || expression == nil {
		return nil, errors.New("invalid parameters")
	}
	return &Return{Base: processors.Base{Walker: walk, Expression: expression}}, nil
}

func (t *Return) Process(node ast.Node, root renders.Node, previous renders.Node) error {
	ret, ok := node.(*ast.ReturnStmt)
	if !ok {
		return fmt.Errorf("The node %v is not a *ast.ReturnStmt", node)
	}
	act := &renders.Return{
		BaseNode: renders.BaseNode{
			Par:  root,
			Prev: previous,
			Dep:  root.Depth() + 1,
		},
		Values: []renders.Expr{},
	}

	for _, exp := range ret.Results {
		act.Values = append(act.Values, t.Expression(exp))
	}

	if previous != nil {
		previous.SetNext(act)
	}
	return nil
}
