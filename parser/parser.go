package parser

import (
	"errors"
	"strings"
)

func ParseString(code string) (*Structogram, error) {
	var s Structogram

	tokens, err := tokenize(code)
	if err != nil {
		return nil, errors.New("parse error: " + err.Error())
	}

	position := func() Pos {
		return Pos{Col: tokens[0].col, Line: tokens[0].line}
	}
	endPosition := func() Pos {
		return Pos{Col: tokens[1].col, Line: tokens[1].line}
	}
	skipSpace := func() {
		for len(tokens) > 0 && tokens[0].typ == tokenSpace {
			tokens = tokens[1:]
		}
	}
	skip := func() {
		tokens = tokens[1:]
		skipSpace()
	}
	sees := func(typ tokenType) bool {
		return len(tokens) > 0 && tokens[0].typ == typ
	}
	seesID := func(id string) bool {
		return sees(tokenID) && tokens[0].text == id
	}
	eatString := func() string {
		if sees(tokenString) {
			s := tokens[0].text
			skip()
			return escapeString(s)
		}
		// TODO Have position info in error messages.
		err = errors.New("parse error: string expected")
		return ""
	}
	eat := func(typ tokenType) {
		if sees(typ) {
			skip()
		} else {
			err = errors.New("parse error: " + typ.String() + " expected")
		}
	}

	var parseStatements func() []Statement

	parseBlock := func() Block {
		var b Block
		b.start = position()
		eat('{')
		b.Statements = parseStatements()
		b.end = endPosition()
		eat('}')
		return b
	}

	parseStatement := func() (Statement, bool) {
		if sees(tokenString) {
			var i Instruction
			i.start = position()
			i.end = endPosition()
			i.quoted = tokens[0].text
			i.Text = eatString()
			return i, true
		} else if seesID("if") {
			ifStart := position()
			skip()
			var condition String
			condition.start = position()
			condition.end = endPosition()
			condition.quoted = tokens[0].text
			condition.Text = eatString()
			then := parseBlock()
			if seesID("else") {
				skip()
				return IfElse{
					start:     ifStart,
					Condition: condition,
					Then:      then,
					Else:      parseBlock(),
				}, true
			}
			return If{
				start:     ifStart,
				Condition: condition,
				Then:      then,
			}, true
		} else if seesID("switch") {
			skip()
			var switchStmt Switch
			switchStmt.Subject = eatString()
			eat('{')
			for seesID("case") {
				skip()
				var c SwitchCase
				if seesID("default") {
					skip()
					c.IsDefault = true
				} else {
					c.Condition = eatString()
				}
				c.Block = parseBlock()
				switchStmt.Cases = append(switchStmt.Cases, c)
			}
			eat('}')
			return switchStmt, true
		} else if seesID("while") {
			skip()
			if sees(tokenString) {
				return While{
					Condition: eatString(),
					Block:     parseBlock(),
				}, true
			} else {
				return InfiniteLoop{Block: parseBlock()}, true
			}
		} else if seesID("do") {
			skip()
			var do DoWhile
			do.Block = parseBlock()
			if seesID("while") {
				skip()
			} else {
				err = errors.New("keyword 'while' expected at the end of do-while loop")
				return nil, false
			}
			do.Condition = eatString()
			return do, true
		} else if seesID("break") {
			var b Break
			b.start = position()
			skip()
			b.quoted = tokens[0].text
			b.end = endPosition()
			b.Text = eatString()
			return b, true
		} else if seesID("call") {
			var c Call
			c.start = position()
			skip()
			c.quoted = tokens[0].text
			c.end = endPosition()
			c.Text = eatString()
			return c, true
		} else if seesID("parallel") {
			skip()
			var p Parallel
			eat('{')
			for sees('{') {
				p.Blocks = append(p.Blocks, parseBlock())
			}
			eat('}')
			return p, true
		}
		return nil, false
	}

	parseStatements = func() []Statement {
		var all []Statement
		for {
			if s, ok := parseStatement(); ok {
				all = append(all, s)
			} else {
				break
			}
		}
		return all
	}

	// The code might start with white space, we want to skip it.
	skipSpace()
	// Parse optional title.
	if seesID("title") {
		skip()
		s.Title.start = position()
		s.Title.quoted = tokens[0].text
		s.Title.Text = eatString()
		s.Title.end = position()
	}
	// Parse code.
	s.Statements = parseStatements()

	// We might have set the err variable and if we have, we do not want to
	// return a half-backed structogram so we return either nil and the error or
	// the strucogram and nil.
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func escapeString(s string) string {
	s = s[1 : len(s)-1] // Trim '"' at front and back.
	s = strings.Replace(s, `\n`, "\n", -1)
	s = strings.Replace(s, `\\`, "\\", -1)
	s = strings.Replace(s, `\"`, "\"", -1)
	return s
}
