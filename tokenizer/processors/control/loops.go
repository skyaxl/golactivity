package control

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type For struct {
	processors.Base
}

// NewSwitchProcessor creates new switch processor to generate renders.If
func NewForProcessor(walk processors.Walker, expression processors.Expression) (processors.Processor, error) {
	if walk == nil || expression == nil {
		return nil, errors.New("invalid parameters")
	}
	return &For{Base: processors.Base{Walker: walk, Expression: expression}}, nil
}

// Process for loop conversion
func (t *For) Process(node ast.Node, root renders.Node, previous renders.Node) error {
	fors, ok := node.(*ast.ForStmt)
	if !ok {
		return fmt.Errorf("The node %v is not a *ast.ForStmt", node)
	}
	forr := &renders.For{
		BaseNode: renders.BaseNode{
			Par:  root,
			Prev: previous,
			Dep:  root.Depth() + 1,
		},
		Conditions: t.Expression(fors.Cond),
	}
	if fors.Init != nil {
		forr.Init = &renders.Assignation{
			BaseNode: renders.BaseNode{
				Par: forr,
				Dep: root.Depth() + 1,
			},
			Left:  make([]renders.Expr, 0),
			Right: make([]renders.Expr, 0),
		}
		if ass, ok := fors.Init.(*ast.AssignStmt); ok {
			for _, ex := range ass.Lhs {
				forr.Init.Left = append(forr.Init.Left, t.Expression(ex))
			}
			for _, ex := range ass.Rhs {
				forr.Init.Right = append(forr.Init.Right, t.Expression(ex))
			}
		}
	}
	//@TODO
	// if fors.Post != nil {

	// }

	forr.Body = &renders.Root{}
	forr.Body.Par = forr

	t.Walker(forr, forr.Body, fors.Body)
	if previous != nil {
		previous.SetNext(forr)
	}
	return nil
}

type Range struct {
	processors.Base
}

// NewRangeProcessor creates new range processor to generate renders.If
func NewRangeProcessor(walk processors.Walker, expression processors.Expression) (processors.Processor, error) {
	if walk == nil || expression == nil {
		return nil, errors.New("invalid parameters")
	}
	return &Range{Base: processors.Base{Walker: walk, Expression: expression}}, nil
}

func (t *Range) Process(node ast.Node, root renders.Node, previous renders.Node) error {
	r, ok := node.(*ast.RangeStmt)
	if !ok {
		return fmt.Errorf("The node %v is not a *ast.RangeStmt", node)
	}
	rangerr := &renders.Range{
		BaseNode: renders.BaseNode{
			Par:  root,
			Prev: previous,
			Dep:  root.Depth() + 1,
		},
		ID:    t.Expression(r.X),
		Key:   t.Expression(r.Key),
		Value: t.Expression(r.Value),
		Body:  &renders.Root{},
	}

	rangerr.Body.Par = rangerr
	t.Walker(rangerr, rangerr.Body, r.Body)
	if previous != nil {
		previous.SetNext(rangerr)
	}
	return nil
}
