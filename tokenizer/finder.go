package tokenizer

import (
	"go/ast"
	"strings"
)

type Aspect string

const (
	DrawAspect Aspect = "@draw"
)

func ExistsComment(comments []*ast.Comment, aspect Aspect) bool {
	for _, comment := range comments {
		if strings.Contains(comment.Text, string(aspect)) {
			return true
		}
	}
	return false
}
