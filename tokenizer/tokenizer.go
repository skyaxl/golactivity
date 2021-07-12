package tokenizer

import (
	"fmt"
	"go/ast"
	"os"
	"reflect"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
	"github.com/skyaxl/golactivity/tokenizer/processors/blocks"
	"github.com/skyaxl/golactivity/tokenizer/processors/control"

	"github.com/withmandala/go-log"
)

var (
	typeProcessors = map[reflect.Type]processors.ProcessorCreator{
		reflect.TypeOf(&ast.IfStmt{}):     control.NewIfProcessor,
		reflect.TypeOf(&ast.SwitchStmt{}): control.NewSwitchProcessor,
		reflect.TypeOf(&ast.ForStmt{}):    control.NewForProcessor,
		reflect.TypeOf(&ast.RangeStmt{}):  control.NewRangeProcessor,
		reflect.TypeOf(&ast.ReturnStmt{}): control.NewReturnProcessor,
		reflect.TypeOf(&ast.BlockStmt{}):  blocks.NewBlockProcessor,
		reflect.TypeOf(&ast.ExprStmt{}):   blocks.NewExprProcessor,
		reflect.TypeOf(&ast.AssignStmt{}): blocks.NewAssignationProcessor,
		reflect.TypeOf(&ast.CallExpr{}):   blocks.NewCallProcessor,
	}
	logger = log.New(os.Stderr)
)

type Transformer struct {
	funcs *ast.FuncDecl
}

//NewTransformer new transformer
func NewTransformer(funcs *ast.FuncDecl) *Transformer {
	return &Transformer{funcs: funcs}
}

// Transform the fun in document
func (t Transformer) Transform() (doc *renders.Document) {
	fun := t.funcs
	doc = &renders.Document{
		Comment: fun.Doc.Text(),
		Name:    fun.Name.Name,
	}
	doc.Root = &renders.Root{}
	doc.Root.Params = FieldListToMap(fun.Type.Params, t)
	if fun.Type.Results != nil {
		doc.Root.Responses = FieldListToMap(fun.Type.Results, t)
	}
	t.Walk(doc.Root, nil, fun.Body)

	return doc
}

func (t Transformer) Walk(root renders.Node, previous renders.Node, node ast.Node) {
	var (
		creator   processors.ProcessorCreator
		ok        bool
		err       error
		processor processors.Processor
	)
	if previous == nil {
		previous = root
	}
	if creator, ok = typeProcessors[reflect.TypeOf(node)]; !ok {
		logger.Errorf("The type %T was not found in processors", node)
		return
	}

	if processor, err = creator(t.Walk, t.Expressions); err != nil {
		logger.Fatalf("A fatal error was found: %v", err)
		return
	}
	if err = processor.Process(node, root, previous); err != nil {
		logger.Fatalf("A fatal error was found to process: %v", err)
		return
	}
	return
}

//Expressions transform expressions
// *ast.Binary, *ast.Unary, *ast.Ident, *ast.BasicLit
func (t Transformer) Expressions(exp ast.Expr) renders.Expr {
	switch exp.(type) {
	case *ast.UnaryExpr:
		{
			u := exp.(*ast.UnaryExpr)
			du := renders.Unary{}
			du.Oper = u.Op.String()
			du.Left = t.Expressions(u.X)
			return du
		}
	case *ast.Ident:
		{
			u := exp.(*ast.Ident)
			du := renders.Identifier{}
			du.ID = u.String()
			return du
		}
	case *ast.BinaryExpr:
		{
			u := exp.(*ast.BinaryExpr)
			du := renders.Binary{}
			du.Oper = u.Op.String()
			du.Left = t.Expressions(u.X)
			du.Right = t.Expressions(u.Y)
			return du
		}
	case *ast.ParenExpr:
		{
			u := exp.(*ast.ParenExpr)
			du := renders.Parent{}
			du.Expr = t.Expressions(u.X)
			return du
		}
	case *ast.BasicLit:
		{
			u := exp.(*ast.BasicLit)
			v := renders.Value{}
			v.Value = u.Value
			v.Kind = u.Kind.String()
			return v
		}
	case *ast.ChanType:
		{
			u := exp.(*ast.ChanType)
			v := renders.Chan{}
			v.Value = t.Expressions(u.Value)
			return v
		}
	case *ast.FuncLit:
		{
			u := exp.(*ast.FuncLit)
			v := renders.FunLiteral{
				Args:      make(renders.Expressions, 0),
				Responses: make(renders.Expressions, 0),
			}
			for _, a := range u.Type.Params.List {
				v.Args = append(v.Args, renders.Field{
					Name: GetName(a.Names),
					Kind: t.Expressions(a.Type),
				})
			}
			return v
		}
	case *ast.SelectorExpr:
		{
			u := exp.(*ast.SelectorExpr)
			v := renders.Identifier{}
			x := t.Expressions(u.X)
			sel := t.Expressions(u.Sel)
			v.ID = fmt.Sprintf("%s.%s", x.String(), sel.String())
			return v
		}
	case *ast.CallExpr:
		{
			u := exp.(*ast.CallExpr)
			v := renders.Call{
				Func:      t.Expressions(u.Fun).(renders.Identifier),
				Arguments: make(renders.Expressions, 0),
			}
			for _, arg := range u.Args {
				v.Arguments = append(v.Arguments, t.Expressions(arg))
			}

			return v
		}
	case *ast.CompositeLit:
		{
			u := exp.(*ast.CompositeLit)
			v := renders.Literal{
				Kind:     t.Expressions(u.Type),
				Elements: make(renders.Expressions, 0),
			}

			for _, arg := range u.Elts {
				v.Elements = append(v.Elements, t.Expressions(arg))
			}
			return v
		}
	case *ast.ArrayType:
		{
			u := exp.(*ast.ArrayType)
			v := renders.ArrayType{
				Type: t.Expressions(u.Elt),
			}
			if u.Len != nil {
				v.Len = t.Expressions(u.Len)
			}
			return v
		}
	case *ast.IndexExpr:
		{
			u := exp.(*ast.IndexExpr)
			v := renders.Index{
				Ident: t.Expressions(u.X),
				Index: t.Expressions(u.Index),
			}
			return v
		}
	}

	return nil
}
