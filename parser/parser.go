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
	Else      Block
}

type Switch struct {
	Subject string
	Cases   []SwitchCase
	Default Block
}

type SwitchCase struct {
	Condition string
	Block     Block
}

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
			var ifElse If
			ifElse.Condition = eatString()
			ifElse.Then = parseBlock()
			if seesID("else") {
				skip()
				ifElse.Else = parseBlock()
			}
			return ifElse, true
		} else if seesID("switch") {
			skip()
			var switchStmt Switch
			switchStmt.Subject = eatString()
			eat('{')
			for seesID("case") {
				skip()
				if seesID("default") {
					skip()
					switchStmt.Default = parseBlock()
				} else {
					condition := eatString()
					block := parseBlock()
					switchStmt.Cases = append(switchStmt.Cases, SwitchCase{
						Condition: condition,
						Block:     block,
					})
				}
			}
			eat('}')
			return switchStmt, true
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
	s = s[1 : len(s)-1]
	s = strings.Replace(s, `\n`, "\n", -1)
	s = strings.Replace(s, `\\`, "\\", -1)
	s = strings.Replace(s, `\"`, "\"", -1)
	return s
}
