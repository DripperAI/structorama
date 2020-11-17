package parser

import (
	"github.com/gonutz/check"

	"testing"
)

func TestTokenization(t *testing.T) {
	tokens, err := tokenize(`title(){}"" "î" "\n\\\""` + "\n\tNextLine")
	check.Eq(t, err, nil)
	tok := func(want token) {
		t.Helper()
		check.Eq(t, tokens[0], want)
		tokens = tokens[1:]
	}

	tok(token{typ: tokenID, text: "title", col: 1, line: 1})
	tok(token{typ: '(', text: "(", col: 6, line: 1})
	tok(token{typ: ')', text: ")", col: 7, line: 1})
	tok(token{typ: '{', text: "{", col: 8, line: 1})
	tok(token{typ: '}', text: "}", col: 9, line: 1})
	tok(token{typ: tokenString, text: `""`, col: 10, line: 1})
	tok(token{typ: tokenSpace, text: " ", col: 12, line: 1})
	tok(token{typ: tokenString, text: `"î"`, col: 13, line: 1})
	tok(token{typ: tokenSpace, text: " ", col: 16, line: 1})
	tok(token{typ: tokenString, text: `"\n\\\""`, col: 17, line: 1})
	tok(token{typ: tokenSpace, text: "\n\t", col: 25, line: 1})
	tok(token{typ: tokenID, text: "NextLine", col: 2, line: 2})

	check.Eq(t, len(tokens), 0) // All tokens checked off the list.
}

func TestTokenizingEscapeSequences(t *testing.T) {
	tokens, err := tokenize(`"quote:\" backslash:\\ line-break:\n"`)
	check.Eq(t, err, nil)
	check.Eq(t, len(tokens), 1)
}

func TestEmptyStringYieldsEmptyStructogram(t *testing.T) {
	s, err := ParseString("")
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{})
}

func TestTitleComesFirst(t *testing.T) {
	s, err := ParseString(`title "the title"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Title: "the title"})
}

func TestTitleStringIsEscaped(t *testing.T) {
	s, err := ParseString(`title "quote:\" backslash:\\ line-break:\n"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Title: "quote:\" backslash:\\ line-break:\n"})
}

func TestRegularInstructionsAreJustStrings(t *testing.T) {
	s, err := ParseString(`"instruction"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Instruction("instruction"),
	}})
}

func TestIfCanHaveNoElse(t *testing.T) {
	s, err := ParseString(`
if "condition" {
	"do this"
	"and that"
}
`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		If{
			Condition: "condition",
			Then: Block{
				Instruction("do this"),
				Instruction("and that"),
			},
		},
	}})
}

func TestIfCanHaveElse(t *testing.T) {
	s, err := ParseString(`
if "false" {
	"then this"
} else {
	"else this"
}
`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		If{
			Condition: "false",
			Then: Block{
				Instruction("then this"),
			},
			Else: Block{
				Instruction("else this"),
			},
		},
	}})
}

func TestSwitchCanBeEmpty(t *testing.T) {
	s, err := ParseString(`
switch "thing" {}
`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Switch{Subject: "thing"},
	}})
}

func TestSwitchWithCases(t *testing.T) {
	s, err := ParseString(`
switch "x" {
	case "1" {}
	case "2" { "two" }
}
`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Switch{
			Subject: "x",
			Cases: []SwitchCase{
				{Condition: "1"},
				{Condition: "2", Block: Block{Instruction("two")}},
			},
		},
	}})
}

func TestSwitchWithDefaultCases(t *testing.T) {
	s, err := ParseString(`
switch "x" {
	case "1" {}
	case default { "default" }
}
`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Switch{
			Subject: "x",
			Cases: []SwitchCase{
				{Condition: "1"},
			},
			Default: Block{Instruction("default")},
		},
	}})
}

func TestInfiniteWhileLoopHasNoCondition(t *testing.T) {
	s, err := ParseString(`while { "do" }`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		While{
			Block: Block{Instruction("do")},
		},
	}})
}

func TestWhileLoopHasConditionNextToTheWhileKeyword(t *testing.T) {
	s, err := ParseString(`while "condition" { "do" }`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		While{
			Condition: "condition",
			Block:     Block{Instruction("do")},
		},
	}})
}
