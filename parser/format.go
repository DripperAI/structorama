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
		if x.Title.quoted != "" {
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
	case If:
		p.WriteString("if ")
		p.WriteString(x.Condition.quoted)
		if x.TrueText.quoted != "" {
			p.WriteString(" ")
			p.WriteString(x.TrueText.quoted)
		}
		p.WriteString(" {")
		p.indentRight()
		p.newLine()
		p.print(x.Then)
		p.indentLeft()
		p.newLine()
		p.WriteString("}")
	case IfElse:
		p.WriteString("if ")
		p.WriteString(x.Condition.quoted)
		if x.TrueText.quoted != "" {
			p.WriteString(" ")
			p.WriteString(x.TrueText.quoted)
		}
		p.WriteString(" {")
		p.indentRight()
		p.newLine()
		p.print(x.Then)
		p.indentLeft()
		p.newLine()
		p.WriteString("} else ")
		if x.FalseText.quoted != "" {
			p.WriteString(x.FalseText.quoted)
			p.WriteString(" ")
		}
		p.WriteString("{")
		p.indentRight()
		p.newLine()
		p.print(x.Else)
		p.indentLeft()
		p.newLine()
		p.WriteString("}")
	case Block:
		for i, stmt := range x.Statements {
			if i > 0 {
				p.newLine()
			}
			p.print(stmt)
		}
	case Switch:
		p.WriteString("switch ")
		p.WriteString(x.Subject.quoted)
		p.WriteString(" {")
		p.indentRight()
		p.newLine()
		for i, c := range x.Cases {
			if i > 0 {
				p.newLine()
			}
			if c.IsDefault {
				p.WriteString("case default {")
			} else {
				p.WriteString("case ")
				p.WriteString(c.Condition.quoted)
				p.WriteString(" {")
			}
			p.indentRight()
			p.newLine()
			p.print(c.Block)
			p.indentLeft()
			p.newLine()
			p.WriteString("}")
		}
		p.indentLeft()
		p.newLine()
		p.WriteString("}")
	case InfiniteLoop:
		p.WriteString("while {")
		p.indentRight()
		p.newLine()
		p.print(x.Block)
		p.indentLeft()
		p.newLine()
		p.WriteString("}")
	case While:
		p.WriteString("while ")
		p.WriteString(x.Condition.quoted)
		p.WriteString(" {")
		p.indentRight()
		p.newLine()
		p.print(x.Block)
		p.indentLeft()
		p.newLine()
		p.WriteString("}")
	case DoWhile:
		p.WriteString("do {")
		p.indentRight()
		p.newLine()
		p.print(x.Block)
		p.indentLeft()
		p.newLine()
		p.WriteString("} while ")
		p.WriteString(x.Condition.quoted)
	case Parallel:
		p.WriteString("parallel {")
		p.indentRight()
		p.newLine()
		for i, b := range x.Blocks {
			if i > 0 {
				p.newLine()
			}
			p.WriteString("{")
			p.indentRight()
			p.newLine()
			p.print(b)
			p.indentLeft()
			p.newLine()
			p.WriteString("}")
		}
		p.indentLeft()
		p.newLine()
		p.WriteString("}")
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
