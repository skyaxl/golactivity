package tokenizer

import (
	"go/ast"
	"os"
	"reflect"

	"github.com/skyaxl/golactivity/renders"
	"github.com/skyaxl/golactivity/tokenizer/processors"
	"github.com/skyaxl/golactivity/tokenizer/processors/blocks"
	"github.com/skyaxl/golactivity/tokenizer/processors/control"
	"github.com/skyaxl/golactivity/tokenizer/processors/expressions"

	"github.com/withmandala/go-log"
)

var (
	statementProcessors = map[reflect.Type]processors.ProcessorCreator{
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
	expressionProcessors = map[reflect.Type]processors.ExpressionProcessorCreator{
		reflect.TypeOf(&ast.UnaryExpr{}):    expressions.NewUnary,
		reflect.TypeOf(&ast.Ident{}):        expressions.NewIdent,
		reflect.TypeOf(&ast.BinaryExpr{}):   expressions.NewBinary,
		reflect.TypeOf(&ast.ParenExpr{}):    expressions.NewParen,
		reflect.TypeOf(&ast.BasicLit{}):     expressions.NewBasicLit,
		reflect.TypeOf(&ast.ChanType{}):     expressions.NewChan,
		reflect.TypeOf(&ast.FuncLit{}):      expressions.NewFuncLit,
		reflect.TypeOf(&ast.SelectorExpr{}): expressions.NewSelector,
		reflect.TypeOf(&ast.CallExpr{}):     expressions.NewCall,
		reflect.TypeOf(&ast.CompositeLit{}): expressions.NewCompositeLit,
		reflect.TypeOf(&ast.ArrayType{}):    expressions.NewArrayType,
		reflect.TypeOf(&ast.IndexExpr{}):    expressions.NewIndex,
	}
	logger = log.New(os.Stderr)
)

// RegisterStatementProcessor add new processor or replace existent
func RegisterStatementProcessor(tp reflect.Type, processor processors.ExpressionProcessorCreator) {
	expressionProcessors[tp] = processor
}

// RegisterExpressionProcessor add new processor or replace existent
func RegisterExpressionProcessor(tp reflect.Type, processor processors.ProcessorCreator) {
	statementProcessors[tp] = processor
}

// Transformer responsible to conver ast tokens to renders tokens.
type Transformer struct {
	funcs *ast.FuncDecl
}

// NewTransformer new transformer
// receives a func declaration
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

// Walk navigates to nodes and process them.
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
	if creator, ok = statementProcessors[reflect.TypeOf(node)]; !ok {
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
}

//Expressions transform expressions
// *ast.Binary, *ast.Unary, *ast.Ident, *ast.BasicLit
func (t Transformer) Expressions(exp ast.Expr) (res renders.Expr) {
	var (
		creator   processors.ExpressionProcessorCreator
		processor processors.ExpressionProcessor
		ok        bool
		err       error
	)
	if creator, ok = expressionProcessors[reflect.TypeOf(exp)]; !ok {
		logger.Errorf("The type %T was not found in processors", exp)
		return nil
	}

	if processor, err = creator(t.Expressions); err != nil {
		logger.Fatalf("A fatal error was found in create a processor: %v", err)
		return nil
	}

	if res, err = processor.Process(exp); err != nil {
		logger.Fatalf("A fatal error was found to process: %v", err)
		return nil
	}

	return res
}
