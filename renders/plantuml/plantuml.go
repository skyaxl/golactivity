package plantuml

import (
	"fmt"
	"io"
	"strings"

	drawer "github.com/skyaxl/golactivity/renders"
)

type Plantuml struct {
	doc *drawer.Document
	wr  io.Writer
}

func NewDrawer(doc *drawer.Document, writter io.Writer) *Plantuml {
	return &Plantuml{doc, writter}
}

func (p *Plantuml) Start() {
	_, _ = io.WriteString(p.wr, "@startuml\nstart\n")
	_, _ = io.WriteString(p.wr, fmt.Sprintf("partition %s {", p.doc.Name))
}

func (p *Plantuml) End() {
	_, _ = io.WriteString(p.wr, "\n}\nstop\n@enduml")
}

func (p *Plantuml) Node(n drawer.Node) {

	if n == nil {
		return
	}

	switch nd := n.(type) {
	case *drawer.If:
		{
			p.If((*drawer.If)(nd))
		}
	case *drawer.Activity:
		{
			p.Activity((*drawer.Activity)(nd))

		}
	case *drawer.While:
		{
			p.While((*drawer.While)(nd))
		}
	case *drawer.Return:
		{
			p.Return((*drawer.Return)(nd))
		}
	case *drawer.For:
		{
			p.For((*drawer.For)(nd))
		}
	case *drawer.Assignation:
		{
			p.Assignation((*drawer.Assignation)(nd))
		}
	case *drawer.Range:
		{
			p.Range((*drawer.Range)(nd))
		}
	case *drawer.Switch:
		{
			p.Switch((*drawer.Switch)(nd))
		}
	}

	if n == nil {
		return
	}

	if n.Next() != nil {
		p.Node(n.Next())
	}
}

func (p *Plantuml) If(i *drawer.If) {
	tab := strings.Repeat("\t", i.Dep)
	init := ""
	if i.Init != nil {
		init = i.Init.String()
		_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%s:%s;", tab, init))
	}
	cond := i.Conditions.String()
	line := fmt.Sprintf("\n%sif (%s) then (yes)", tab, cond)
	_, _ = io.WriteString(p.wr, line)
	p.Node(i.Body)
	_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%selse", tab))
	if i.Else != nil {
		p.Node(i.Else)
	}
	_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%sendif", tab))
}

func (p *Plantuml) While(i *drawer.While) {
	tab := strings.Repeat("\t", i.Dep)
	_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%srepeate (%s)", tab, i.Conditions.String()))
	p.Node(i.Body)
}

func (p *Plantuml) Activity(a *drawer.Activity) {
	tab := strings.Repeat("\t", a.Dep)
	_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%s:%s;", tab, a.Exp))
	if a.Comment != "" {
		_, _ = io.WriteString(p.wr, fmt.Sprintf("%snote right:%s;", tab, a.Comment))
	}
}

func (p *Plantuml) Return(a *drawer.Return) {
	tab := strings.Repeat("\t", a.Dep)
	_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%send", tab))
	_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%snote right:Return (%s)", tab, a.Values.Join(",")))
}

func (p *Plantuml) For(i *drawer.For) {
	tab := strings.Repeat("\t", i.Dep)
	init := ""
	if i.Init != nil {
		init = i.Init.String()
		_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%s:%s;", tab, init))
	}
	cond := i.Conditions.String()
	line := fmt.Sprintf("\n%srepeat", tab)
	_, _ = io.WriteString(p.wr, line)
	p.Node(i.Body)
	line = fmt.Sprintf("\n%srepeat while (%s) is (true)", tab, cond)
	_, _ = io.WriteString(p.wr, line)
}

func (p *Plantuml) Range(r *drawer.Range) {
	tab := strings.Repeat("\t", r.Dep)

	line := fmt.Sprintf("\n%srepeat", tab)
	_, _ = io.WriteString(p.wr, line)
	key := r.Key.String()
	if key != "_" {
		_, _ = io.WriteString(p.wr, fmt.Sprintf("\n\t%s:%s := keyOf %s;", tab, key, r.ID))
	}
	if r.Value != nil && r.Value.String() != "_" {
		_, _ = io.WriteString(p.wr, fmt.Sprintf("\n\t%s:%s := itemOf %s;", tab, r.Value, r.ID))
	}
	p.Node(r.Body)
	line = fmt.Sprintf("\n%srepeat while (range %s) is (true)", tab, r.ID)
	_, _ = io.WriteString(p.wr, line)
}

func (p *Plantuml) Assignation(a *drawer.Assignation) {
	tab := strings.Repeat("\t", a.Dep)
	_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%s:%s;", tab, a.String()))
}

func (p *Plantuml) Switch(s *drawer.Switch) {
	tab := strings.Repeat("\t", s.Dep)
	init := ""
	if s.Init != nil {
		init = s.Init.String()
		_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%s:%s;", tab, init))
	}

	cond := s.Tag.String()
	line := fmt.Sprintf("\n\n%sswitch (switch(%s))", tab, cond)
	_, _ = io.WriteString(p.wr, line)
	for _, c := range s.Cases {
		if c.Value != nil && len(c.Value) != 0 {
			_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%scase (case %s)", tab, c.Value.Join("&&")))
		} else {
			_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%scase (default)", tab))
		}
		p.Node(c.Body)
		_, _ = io.WriteString(p.wr, "\n")
	}
	_, _ = io.WriteString(p.wr, fmt.Sprintf("\n%sendswitch", tab))
}
