package plantuml

import (
	"fmt"
	"io"
	"strings"

	"github.com/skyaxl/golactivity/drawer"
)

type Plantuml struct {
	doc *drawer.Document
	wr  io.Writer
}

func NewDrawer(doc *drawer.Document, writter io.Writer) *Plantuml {
	return &Plantuml{doc, writter}
}

func (p *Plantuml) Start() {
	io.WriteString(p.wr, "@startuml\nstart\n")
	io.WriteString(p.wr, fmt.Sprintf("partition %s {", p.doc.Name))
}

func (p *Plantuml) End() {
	io.WriteString(p.wr, "\n}\nstop\n@enduml")
}

func (p *Plantuml) Node(n drawer.Node) {

	if n == nil {
		return
	}

	switch n.(type) {
	case *drawer.If:
		{
			p.If(n.(*drawer.If))
		}
	case *drawer.Activity:
		{
			p.Activity(n.(*drawer.Activity))

		}
	case *drawer.While:
		{
			p.While(n.(*drawer.While))
		}
	case *drawer.Return:
		{
			p.Return(n.(*drawer.Return))
		}
	case *drawer.For:
		{
			p.For(n.(*drawer.For))
		}

	}

	if n == nil {
		return
	}

	if n.Next() != nil {
		p.Node(n.Next())
	}
	return
}

func (p *Plantuml) If(i *drawer.If) {
	tab := strings.Repeat("\t", i.Dep)
	init := ""
	if i.Init != nil {
		init = i.Init.String()
		io.WriteString(p.wr, fmt.Sprintf("\n%s:%s;", tab, init))
	}
	cond := i.Conditions.String()
	line := fmt.Sprintf("\n%sif (%s) then (yes)", tab, cond)
	io.WriteString(p.wr, line)
	p.Node(i.Body)
	io.WriteString(p.wr, fmt.Sprintf("\n%selse", tab))
	if i.Else != nil {
		p.Node(i.Else)
	}
	io.WriteString(p.wr, fmt.Sprintf("\n%sendif", tab))
}

func (p *Plantuml) While(i *drawer.While) {
	tab := strings.Repeat("\t", i.Dep)
	io.WriteString(p.wr, fmt.Sprintf("\n%srepeate (%s)", tab, i.Conditions.String()))
	p.Node(i.Body)
}

func (p *Plantuml) Activity(a *drawer.Activity) {
	tab := strings.Repeat("\t", a.Dep)
	io.WriteString(p.wr, fmt.Sprintf("\n%s:%s;", tab, a.Name))
	if a.Comment != "" {
		io.WriteString(p.wr, fmt.Sprintf("%snote right:%s;", tab, a.Comment))
	}
}

func (p *Plantuml) Return(a *drawer.Return) {
	tab := strings.Repeat("\t", a.Dep)
	io.WriteString(p.wr, fmt.Sprintf("\n%send", tab))
	io.WriteString(p.wr, fmt.Sprintf("\n%snote right:Return (%s)", tab, a.Values.Join(",")))
}

func (p *Plantuml) For(i *drawer.For) {
	tab := strings.Repeat("\t", i.Dep)
	init := ""
	if i.Init != nil {
		init = i.Init.String()
		io.WriteString(p.wr, fmt.Sprintf("\n%s:%s;", tab, init))
	}
	cond := i.Conditions.String()
	line := fmt.Sprintf("\n%srepeat", tab)
	io.WriteString(p.wr, line)
	p.Node(i.Body)
	line = fmt.Sprintf("\n%srepeat while (%s) is (true)", tab, cond)
	io.WriteString(p.wr, line)
}