package parser

import (
	"errors"
	"strings"
)

type Structogram struct {
	Title      string
	Statements []Statement
}

type Block []Statement

type Statement interface {
}

type Instruction string

type If struct {
	Condition string
	Then      Block
}

type IfElse struct {
	Condition string
	Then      Block
	Else      Block
}

type Switch struct {
	Subject string
	Cases   []SwitchCase
}

type SwitchCase struct {
	// IsDefault and Condition are exclusive, a switch case is either the
	// default (case default {...}) or has a Condition (case "condition" {...}).
	IsDefault bool
	Condition string
	Block     Block
}

type Call string

type Parallel []Block

type InfiniteLoop Block

type While struct {
	Condition string
	Block     Block
}

type DoWhile struct {
	Block     Block
	Condition string
}

type Break string

func ParseString(code string) (*Structogram, error) {
	var s Structogram

	tokens, err := tokenize(code)
	if err != nil {
		return nil, errors.New("parse error: " + err.Error())
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
		eat('{')
		s := parseStatements()
		eat('}')
		return Block(s)
	}

	parseStatement := func() (Statement, bool) {
		if sees(tokenString) {
			return Instruction(eatString()), true
		} else if seesID("if") {
			skip()
			condition := eatString()
			then := parseBlock()
			if seesID("else") {
				skip()
				return IfElse{
					Condition: condition,
					Then:      then,
					Else:      parseBlock(),
				}, true
			}
			return If{
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
				return InfiniteLoop(parseBlock()), true
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
			skip()
			return Break(eatString()), true
		} else if seesID("call") {
			skip()
			return Call(eatString()), true
		} else if seesID("parallel") {
			skip()
			var p Parallel
			eat('{')
			for sees('{') {
				p = append(p, parseBlock())
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
		s.Title = eatString()
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
