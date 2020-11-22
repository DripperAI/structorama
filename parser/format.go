package parser

import (
	"bytes"
	"fmt"
	"strings"
)

func FormatString(code string) (string, error) {
	s, err := ParseString(code)
	if err != nil {
		return "", err
	}
	p := &printer{}
	p.print(s)
	if p.err != nil {
		return "", p.err
	}
	return p.String(), nil
}

type printer struct {
	bytes.Buffer
	tabs string
	err  error
}

func (p *printer) String() string {
	return string(p.Bytes())
}

func (p *printer) print(node interface{}) {
	if p.err != nil {
		return
	}
	switch x := node.(type) {
	case *Structogram:
		if x.Title.Text != "" {
			p.WriteString("title ")
			p.WriteString(x.Title.quoted)
			if len(x.Statements) > 0 {
				p.WriteString("\n\n")
			}
		}
		for i, stmt := range x.Statements {
			p.print(stmt)
			if i+1 < len(x.Statements) {
				a := x.Statements[i]
				b := x.Statements[i+1]
				if b.Start().Line-a.End().Line >= 2 {
					p.WriteString("\n")
				}
			}
			p.newLine()
		}
	case Instruction:
		p.WriteString(x.quoted)
	case Call:
		p.WriteString("call ")
		p.WriteString(x.quoted)
	case Break:
		p.WriteString("break ")
		p.WriteString(x.quoted)
	case IfElse:
		p.WriteString("if " + x.Condition.quoted + " {")
		p.indentRight()
		p.newLine()
		p.print(x.Then)
		p.indentLeft()
		p.newLine()
		p.WriteString("} else {")
		p.indentRight()
		p.newLine()
		p.print(x.Else)
		p.indentLeft()
		p.newLine()
		p.WriteString("}")
	case Block:
		for i, stmt := range x {
			if i > 0 {
				p.newLine()
			}
			p.print(stmt)
		}
	default:
		p.err = fmt.Errorf("printer.print: unhandled node type %T", node)
	}
}

func (p *printer) indentRight() {
	p.tabs += "\t"
}

func (p *printer) indentLeft() {
	p.tabs = strings.TrimSuffix(p.tabs, "\t")
}

func (p *printer) newLine() {
	p.WriteString("\n" + p.tabs)
}
