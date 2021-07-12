package blocks

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
)

type Block struct {
	processors.Base
}

// NewBlockProcessor creates new block processor to generate renders.Body or Root
func NewBlockProcessor(walk processors.Walker, expression processors.Expression) (processors.Processor, error) {
	if walk == nil || expression == nil {
		return nil, errors.New("invalid parameters")
	}
	return &Block{Base: processors.Base{Walker: walk, Expression: expression}}, nil
}

func (t *Block) Process(node ast.Node, root renders.Node, previous renders.Node) error {
	b, ok := node.(*ast.BlockStmt)
	if !ok {
		return fmt.Errorf("The node %v is not a *ast.BlockStmt", node)
	}
	for _, n := range b.List {
		t.Walker(root, previous, n)
		if previous != nil && previous.Next() != nil {
			previous = previous.Next()
		}
	}
	return nil
}
