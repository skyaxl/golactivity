package control

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type If struct {
	processors.Base
}

// NewIfProcessor creates new if processor to generate renders.If
func NewIfProcessor(walk processors.Walker, expression processors.Expression) (processors.Processor, error) {
	if walk == nil || expression == nil {
		return nil, errors.New("invalid parameters")
	}
	return &If{Base: processors.Base{Walker: walk, Expression: expression}}, nil
}

// Process
func (i *If) Process(node ast.Node, root renders.Node, previous renders.Node) error {
	fi, ok := node.(*ast.IfStmt)
	if !ok {
		return fmt.Errorf("The node %v is not a *ast.IfStmt", node)
	}
	fiNo := &renders.If{
		BaseNode: renders.BaseNode{
			Par:  root,
			Prev: previous,
			Dep:  root.Depth() + 1,
		},
		Conditions: i.Expression(fi.Cond),
	}
	if fi.Init != nil {
		fiNo.Init = &renders.Assignation{
			BaseNode: renders.BaseNode{
				Par: fiNo,
				Dep: root.Depth() + 1,
			},
			Left:  make([]renders.Expr, 0),
			Right: make([]renders.Expr, 0),
		}
		if ass, ok := fi.Init.(*ast.AssignStmt); ok {
			for _, ex := range ass.Lhs {
				fiNo.Init.Left = append(fiNo.Init.Left, i.Expression(ex))
			}
			for _, ex := range ass.Rhs {
				fiNo.Init.Right = append(fiNo.Init.Right, i.Expression(ex))
			}
		}
	}

	fiNo.Body = &renders.Root{}
	fiNo.Body.Par = fiNo

	i.Walker(fiNo, fiNo.Body, fi.Body)
	if previous != nil {
		previous.SetNext(fiNo)
	}
	return nil
}

type Switch struct {
	processors.Base
}

// NewSwitchProcessor creates new switch processor to generate renders.If
func NewSwitchProcessor(walk processors.Walker, expression processors.Expression) (processors.Processor, error) {
	if walk == nil || expression == nil {
		return nil, errors.New("invalid parameters")
	}
	return &Switch{Base: processors.Base{Walker: walk, Expression: expression}}, nil
}

// Process switch conversion
func (t *Switch) Process(node ast.Node, root renders.Node, previous renders.Node) error {
	sw, ok := node.(*ast.SwitchStmt)
	if !ok {
		return fmt.Errorf("The node %v is not a *ast.SwitchStmt", node)
	}
	swi := &renders.Switch{
		BaseNode: renders.BaseNode{
			Par:  root,
			Prev: previous,
			Dep:  root.Depth() + 1,
		},
		Cases: make([]*renders.Case, 0),
		Tag:   t.Expression(sw.Tag),
	}
	if sw.Init != nil {
		swi.Init = &renders.Assignation{
			BaseNode: renders.BaseNode{
				Par: swi,
				Dep: root.Depth(),
			},
			Left:  make([]renders.Expr, 0),
			Right: make([]renders.Expr, 0),
		}
		if ass, ok := sw.Init.(*ast.AssignStmt); ok {
			for _, ex := range ass.Lhs {
				swi.Init.Left = append(swi.Init.Left, t.Expression(ex))
			}
			for _, ex := range ass.Rhs {
				swi.Init.Right = append(swi.Init.Right, t.Expression(ex))
			}
		}
	}

	if sw.Body != nil {
		for _, c := range sw.Body.List {
			cas := c.(*ast.CaseClause)
			swic := &renders.Case{
				BaseNode: renders.BaseNode{
					Par: swi,
					Dep: swi.Depth(),
				},
				Value: make(renders.Expressions, 0),
				Body:  &renders.Root{},
			}
			swic.Body.BaseNode = renders.BaseNode{
				Par: swic,
				Dep: swic.Depth(),
			}
			if cas.List != nil {
				for _, exp := range cas.List {
					swic.Value = append(swic.Value, t.Expression(exp))
				}
			}

			t.Walker(swic.Body, nil, &ast.BlockStmt{
				List: cas.Body,
			})
			swi.Cases = append(swi.Cases, swic)
		}

	}

	if previous != nil {
		previous.SetNext(swi)
	}

	return nil
}
