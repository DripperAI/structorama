package parser

import (
	"bytes"
	"fmt"
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
	lastLine int
	err      error
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
				p.WriteString("\n")
			}
			// We set the lastLine to -2 so that there is always an empty line
			// between the title and what follows.
			p.lastLine = -2
		} else {
			// Set the lastLine to infinity (kind of) to make sure the first
			// statement will start at the very top.
			p.lastLine = 9999999
		}
		for _, stmt := range x.Statements {
			p.print(stmt)
		}
	case Instruction:
		p.insertEmptyLine(x.span)
		p.WriteString(x.quoted)
		p.WriteString("\n")
	case Call:
		p.insertEmptyLine(x.span)
		p.WriteString("call ")
		p.WriteString(x.quoted)
		p.WriteString("\n")
	case Break:
		p.insertEmptyLine(x.span)
		p.WriteString("break ")
		p.WriteString(x.quoted)
		p.WriteString("\n")
	default:
		p.err = fmt.Errorf("printer.print: unhandled node type %T", node)
	}
}

func (p *printer) insertEmptyLine(at span) {
	if at.start.line-p.lastLine >= 2 {
		p.WriteString("\n")
	}
	p.lastLine = at.end.line
}
