package drawer

import "fmt"

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
	Name    string
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

type Value struct {
	Value string
	Kind  string
}

func (ex Value) String() string {
	return fmt.Sprintf("%s", ex.Value)
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

//If self
type If struct {
	BaseNode
	Conditions Expr
	Body       *Root
	Else       *Root
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
	Values []Expr
}

type Drawer interface {
	Start()
	End()
	Node(n Node)
}
