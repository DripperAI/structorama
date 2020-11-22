package parser

import (
	"testing"

	"github.com/gonutz/check"
)

func TestTokenization(t *testing.T) {
	tokens, err := tokenize(`title{}"" "î" "\n\\\""` + "\n\tNextLine")
	check.Eq(t, err, nil)
	tok := func(want token) {
		t.Helper()
		check.Eq(t, tokens[0], want)
		tokens = tokens[1:]
	}

	tok(token{typ: tokenID, text: "title", col: 1, line: 1})
	tok(token{typ: '{', text: "{", col: 6, line: 1})
	tok(token{typ: '}', text: "}", col: 7, line: 1})
	tok(token{typ: tokenString, text: `""`, col: 8, line: 1})
	tok(token{typ: tokenSpace, text: " ", col: 10, line: 1})
	tok(token{typ: tokenString, text: `"î"`, col: 11, line: 1})
	tok(token{typ: tokenSpace, text: " ", col: 14, line: 1})
	tok(token{typ: tokenString, text: `"\n\\\""`, col: 15, line: 1})
	tok(token{typ: tokenSpace, text: "\n\t", col: 23, line: 1})
	tok(token{typ: tokenID, text: "NextLine", col: 2, line: 2})
	tok(token{typ: tokenEOF, text: "", col: 10, line: 2})

	check.Eq(t, len(tokens), 0) // All tokens checked off the list.
}

func TestTokenizingEscapeSequences(t *testing.T) {
	tokens, err := tokenize(`"quote:\" backslash:\\ line-break:\n"`)
	check.Eq(t, err, nil)
	check.Eq(t, len(tokens), 2)
	check.Eq(t, tokens[0].typ, tokenString)
	check.Eq(t, tokens[0].text, `"quote:\" backslash:\\ line-break:\n"`)
	check.Eq(t, tokens[1].typ, tokenEOF)
}

func TestEmptyStringYieldsEmptyStructogram(t *testing.T) {
	s, err := ParseString("")
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{})
}

func TestTitleComesFirst(t *testing.T) {
	s, err := ParseString(`title "the title"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Title: String{
		Text:   "the title",
		quoted: `"the title"`,
		span: span{
			start: pos{col: 7, line: 1},
			end:   pos{col: 18, line: 1},
		},
	}})
}

func TestTitleStringIsEscaped(t *testing.T) {
	s, err := ParseString(`title "quote:\" backslash:\\ line-break:\n"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Title: String{
		Text:   "quote:\" backslash:\\ line-break:\n",
		quoted: `"quote:\" backslash:\\ line-break:\n"`,
		span: span{
			start: pos{col: 7, line: 1},
			end:   pos{col: 44, line: 1},
		},
	}})
}

func TestRegularInstructionsAreJustStrings(t *testing.T) {
	s, err := ParseString(`"instruction"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Instruction{
			Text:   "instruction",
			quoted: `"instruction"`,
			span: span{
				start: pos{col: 1, line: 1},
				end:   pos{col: 14, line: 1},
			},
		},
	}})
}

func TestIfHasNoElse(t *testing.T) {
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
				Instruction{
					Text:   "do this",
					quoted: `"do this"`,
					span: span{
						start: pos{col: 2, line: 3},
						end:   pos{col: 11, line: 3},
					},
				},
				Instruction{
					Text:   "and that",
					quoted: `"and that"`,
					span: span{
						start: pos{col: 2, line: 4},
						end:   pos{col: 12, line: 4},
					},
				},
			},
		},
	}})
}

func TestIfElseHasBoth(t *testing.T) {
	s, err := ParseString(`
if "false" {
	"then this"
} else {
	"else this"
}
`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		IfElse{
			Condition: "false",
			Then: Block{
				Instruction{
					Text:   "then this",
					quoted: `"then this"`,
					span: span{
						start: pos{col: 2, line: 3},
						end:   pos{col: 13, line: 3},
					},
				},
			},
			Else: Block{
				Instruction{
					Text:   "else this",
					quoted: `"else this"`,
					span: span{
						start: pos{col: 2, line: 5},
						end:   pos{col: 13, line: 5},
					},
				},
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
				{Condition: "2", Block: Block{Instruction{
					Text:   "two",
					quoted: `"two"`,
					span: span{
						start: pos{col: 13, line: 4},
						end:   pos{col: 18, line: 4},
					},
				}}},
			},
		},
	}})
}

func TestSwitchWithDefaultCases(t *testing.T) {
	s, err := ParseString(`
switch "x" {
	case "1" {}
	case default {}
}
`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Switch{
			Subject: "x",
			Cases: []SwitchCase{
				{Condition: "1"},
				{IsDefault: true},
			},
		},
	}})
}

func TestInfiniteLoopHasNoCondition(t *testing.T) {
	s, err := ParseString(`while { "do" }`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		InfiniteLoop{Instruction{
			Text:   "do",
			quoted: `"do"`,
			span: span{
				start: pos{col: 9, line: 1},
				end:   pos{col: 13, line: 1},
			},
		}},
	}})
}

func TestWhileLoopHasConditionNextToTheWhileKeyword(t *testing.T) {
	s, err := ParseString(`while "condition" { "do" }`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		While{
			Condition: "condition",
			Block: Block{Instruction{
				Text:   "do",
				quoted: `"do"`,
				span: span{
					start: pos{col: 21, line: 1},
					end:   pos{col: 25, line: 1},
				},
			}},
		},
	}})
}

func TestDoWhileLoopHasConditionInFooter(t *testing.T) {
	s, err := ParseString(`do { "do" } while "condition"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		DoWhile{
			Block: Block{Instruction{
				Text:   "do",
				quoted: `"do"`,
				span: span{
					start: pos{col: 6, line: 1},
					end:   pos{col: 10, line: 1},
				},
			}},
			Condition: "condition",
		},
	}})
}

func TestLoopsCanHaveBreaks(t *testing.T) {
	s, err := ParseString(`while { break "destination" }`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		InfiniteLoop{Break{
			Text:   "destination",
			quoted: `"destination"`,
			span: span{
				start: pos{col: 9, line: 1},
				end:   pos{col: 28, line: 1},
			},
		}},
	}})
}

func TestCallBlockHasOneStringInstruction(t *testing.T) {
	s, err := ParseString(`call "instruction"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Call{
			Text:   "instruction",
			quoted: `"instruction"`,
			span: span{
				start: pos{col: 1, line: 1},
				end:   pos{col: 19, line: 1},
			},
		},
	}})
}

func TestParallelExecutionHasSubBlocks(t *testing.T) {
	s, err := ParseString(`
parallel {
	{
		"block 1"
	}
	{}
	{
		"block 3"
	}
}

parallel {}
`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Parallel{
			Block{Instruction{
				Text:   "block 1",
				quoted: `"block 1"`,
				span: span{
					start: pos{col: 3, line: 4},
					end:   pos{col: 12, line: 4},
				},
			}},
			Block{},
			Block{Instruction{
				Text:   "block 3",
				quoted: `"block 3"`,
				span: span{
					start: pos{col: 3, line: 8},
					end:   pos{col: 12, line: 8},
				},
			}},
		},
		Parallel{},
	}})
}
