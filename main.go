package main

import (
	"os"

	"github.com/skyaxl/golactivity/finder"
	"github.com/skyaxl/golactivity/reader"
	"github.com/skyaxl/golactivity/renders/plantuml"
	"github.com/skyaxl/golactivity/tokenizer"
)

func main() {
	r := reader.New("example", ".*")
	pkgs, fset, _ := r.Read()
	c, _ := os.Getwd()
	println(c)
	funs := finder.FindAllowedFuncs(pkgs, fset)

	for _, fun := range funs {
		t := tokenizer.NewTransformer(fun)
		doc := t.Transform()
		_ = os.Remove("plantum.pl")
		f, _ := os.OpenFile("plantum.pl", os.O_RDWR|os.O_CREATE, os.ModePerm)
		drawer := plantuml.NewDrawer(doc, f)
		drawer.Start()
		drawer.Node(doc.Root)
		drawer.End()
		_ = f.Close()
	}

	//
}
