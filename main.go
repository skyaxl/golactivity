package main

import (
	"os"

	"github.com/skyaxl/golactivity/reader"
	"github.com/skyaxl/golactivity/renders/plantuml"
	"github.com/skyaxl/golactivity/tokenizer"
)

func main() {
	r := reader.New("example", ".*")
	pkgs, fset, _ := r.Read()
	c, _ := os.Getwd()
	println(c)
	funs := tokenizer.ReadTokens(pkgs, fset)

	for _, fun := range funs {
		t := tokenizer.NewTransformer(fun)
		doc := t.Transform()
		drawer := plantuml.NewDrawer(doc, os.Stdout)
		drawer.Start()
		drawer.Node(doc.Root)
		drawer.End()
	}

	//
}
