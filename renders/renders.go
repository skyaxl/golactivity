package renders

import (
	"fmt"
	"strings"
)

//Document doc
type Document struct {
	Name    string
	Comment string
	Root    *Root
}

type Node interface {
	Parent() Node
	Previous() Node
	SetPrevious(Node)

	Next() Node
	SetNext(Node)
	Depth() int
}

type BaseNode struct {
	Par  Node
	Nex  Node
	Prev Node
	Dep  int
}

func (bn *BaseNode) Parent() Node {
	return bn.Par
}

func (bn *BaseNode) Previous() Node {
	return bn.Prev
}

func (bn *BaseNode) Next() Node {
	return bn.Nex
}

func (bn *BaseNode) Depth() int {
	return bn.Dep
}

func (bn *BaseNode) SetNext(n Node) {
	bn.Nex = (n)
}

func (bn *BaseNode) SetPrevious(n Node) {
	bn.Prev = n
}

type Root struct {
	BaseNode
	Params    map[string]string
	Responses map[string]string
}

type Activity struct {
	BaseNode
	Exp     Expr
	Comment string
}

//Repeat todo
type Repeat struct {
	BaseNode
	Conditions Expr
	Comment    string
	Body       []Node
	Backward   Node
}

//Expr need have to string method
type Expr interface {
	String() string
}

//Parent only encapsulate any expression
type Parent struct {
	Expr Expr
}

func (ex Parent) String() string {
	return fmt.Sprintf("(%s)", ex.Expr.String())
}

//Identifier its identifier node
type Identifier struct {
	ID string
}

func (ex Identifier) String() string {
	return ex.ID
}

type Expressions []Expr

func (exps Expressions) Join(sep string) (res string) {
	if exps == nil {
		return res
	}
	for _, ex := range exps {
		if ex == nil {
			continue
		}
		res += fmt.Sprintf("%s,", ex.String())
	}
	if len(res) == 0 {
		return ""
	}
	return string(res[:len(res)-1])
}

//Call its identifier node
type Call struct {
	Func      Expr
	Arguments Expressions
}

func (ex Call) String() string {
	res := fmt.Sprintf("%s(%s)", ex.Func.String(), ex.Arguments.Join(","))
	return res
}

type Value struct {
	Value string
	Kind  string
}

func (ex Value) String() string {
	if ex.Kind == "STRING" {
		return ex.Value
	}
	return fmt.Sprintf("%s(%s)", ex.Kind, ex.Value)
}

type Field struct {
	Name string
	Kind Expr
}

func (ex Field) String() string {
	return fmt.Sprintf("%s %s", ex.Name, ex.Kind)
}

type Literal struct {
	Kind     Expr
	Elements Expressions
}

func (ex Literal) String() string {
	if ex.Kind == nil {
		return fmt.Sprintf("{%s}", ex.Elements.Join(" "))
	}
	if ex.Elements == nil {
		return fmt.Sprintf("%s", ex.Kind.String())
	}

	return fmt.Sprintf("%s(%s)", ex.Kind.String(), ex.Elements.Join(" "))
}

type FunLiteral struct {
	Args      Expressions
	Responses Expressions
}

func (ex FunLiteral) String() string {
	if ex.Responses == nil || len(ex.Responses) == 0 {
		return fmt.Sprintf("func (%s){Literal func}", ex.Args.Join(" "))
	}

	return fmt.Sprintf("func (%s) (%s){Literal func}", ex.Args.Join(","), ex.Responses.Join(","))
}

type Operation struct {
	Oper string
}

//Operator return operator
func (ex Operation) String() string {
	return ex.Oper
}

//Unary if or else
type Unary struct {
	BaseNode
	Operation
	Left    Expr
	Comment string
}

func (ex Unary) String() string {
	return fmt.Sprintf("%s%s", ex.Oper, ex.Left.String())
}

//Binary if or else
type Binary struct {
	BaseNode
	Operation
	Unary
	Right Expr
}

func (ex Binary) String() string {
	return fmt.Sprintf("%s %s %s", ex.Left.String(), ex.Oper, ex.Right.String())
}

type Chan struct {
	Value   Expr
	Comment string
}

func (ex Chan) String() string {
	return fmt.Sprintf("chan %s", ex.Value)
}

//If self
type If struct {
	BaseNode
	Init       *Assignation
	Conditions Expr
	Body       *Root
	Else       *Root
}

//If self
type For struct {
	BaseNode
	Init       *Assignation
	Conditions Expr
	Post       Expr
	Body       *Root
}

//While while representation
type While struct {
	BaseNode
	Conditions Expr
	Body       *Root
}

//Return
type Return struct {
	BaseNode
	Values Expressions
}

//Assignation
type Assignation struct {
	BaseNode
	Left  Expressions
	Right Expressions
}

func (ex Assignation) String() string {
	return fmt.Sprintf("%s := %s", ex.Left.Join(""), ex.Right.Join(""))
}

type ArrayType struct {
	Type Expr
	Len  Expr
}

func (ex ArrayType) String() string {
	l := ""
	if ex.Len != nil {
		l = ex.Len.String()
	}
	return strings.ReplaceAll(fmt.Sprintf("[%s]%s", l, ex.Type), " ", "")
}

type Index struct {
	Ident Expr
	Index Expr
}

func (ex Index) String() string {
	return fmt.Sprintf("%s[%s]", ex.Ident, ex.Index)
}

type Drawer interface {
	Start()
	End()
	Node(n Node)
}

type Range struct {
	BaseNode
	ID    Expr
	Key   Expr
	Value Expr
	Body  *Root
}

type Switch struct {
	BaseNode
	Init  *Assignation
	Tag   Expr
	Cases []*Case
}

type Case struct {
	BaseNode
	Body  *Root
	Value Expressions
}
