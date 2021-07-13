package reader

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
)

// Read an directory and try to find packages
type Reader struct {
	path  string
	regex string
}

//New reader
func New(path string, filterRegex string) *Reader {
	return &Reader{path, filterRegex}
}

func (r *Reader) Read() (res map[string]*ast.Package, fset *token.FileSet, err error) {
	fset = token.NewFileSet()
	res, err = parser.ParseDir(fset, r.path, func(o os.FileInfo) bool {
		if reg, e := regexp.Compile(r.regex); e != nil {
			return reg.MatchString(o.Name())
		}
		return true
	}, parser.ParseComments)
	return res, fset, err
}
