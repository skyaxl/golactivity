package plantuml

import (
	"fmt"
	"io"
	"strings"

	"github.com/skyaxl/golactivity/pkg/drawer"
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
	io.WriteString(p.wr, fmt.Sprintf("\n%sif (%s) then (yes)", tab, i.Conditions.String()))
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
	io.WriteString(p.wr, fmt.Sprintf("\n%snote right:Return (", tab))
	if len(a.Values) != 0 {

		for _, e := range a.Values {
			io.WriteString(p.wr, fmt.Sprintf("%s,", e.String()))
		}
		io.WriteString(p.wr, ");")
	}
}
