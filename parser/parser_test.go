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
	s, err := ParseString(`title "the title" `)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Title: String{
		Text:   "the title",
		quoted: `"the title"`,
		start:  Pos{Col: 7, Line: 1},
		end:    Pos{Col: 18, Line: 1},
	}})
}

func TestTitleStringIsEscaped(t *testing.T) {
	s, err := ParseString(`title "quote:\" backslash:\\ line-break:\n"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Title: String{
		Text:   "quote:\" backslash:\\ line-break:\n",
		quoted: `"quote:\" backslash:\\ line-break:\n"`,
		start:  Pos{Col: 7, Line: 1},
		end:    Pos{Col: 44, Line: 1},
	}})
}

func TestRegularInstructionsAreJustStrings(t *testing.T) {
	s, err := ParseString(`"instruction"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Instruction{
			Text:   "instruction",
			quoted: `"instruction"`,
			start:  Pos{Col: 1, Line: 1},
			end:    Pos{Col: 14, Line: 1},
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
			start: Pos{Col: 1, Line: 2},
			Condition: String{
				Text:   "condition",
				quoted: `"condition"`,
				start:  Pos{Col: 4, Line: 2},
				end:    Pos{Col: 15, Line: 2},
			},
			Then: Block{
				start: Pos{Col: 16, Line: 2},
				end:   Pos{Col: 2, Line: 5},
				Statements: []Statement{
					Instruction{
						Text:   "do this",
						quoted: `"do this"`,
						start:  Pos{Col: 2, Line: 3},
						end:    Pos{Col: 11, Line: 3},
					},
					Instruction{
						Text:   "and that",
						quoted: `"and that"`,
						start:  Pos{Col: 2, Line: 4},
						end:    Pos{Col: 12, Line: 4},
					},
				},
			},
		},
	}})
	check.Eq(t, s.Statements[0].End(), Pos{Col: 2, Line: 5})
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
			start: Pos{Col: 1, Line: 2},
			Condition: String{
				Text:   "false",
				quoted: `"false"`,
				start:  Pos{Col: 4, Line: 2},
				end:    Pos{Col: 11, Line: 2},
			},
			Then: Block{
				start: Pos{Col: 12, Line: 2},
				end:   Pos{Col: 2, Line: 4},
				Statements: []Statement{
					Instruction{
						Text:   "then this",
						quoted: `"then this"`,
						start:  Pos{Col: 2, Line: 3},
						end:    Pos{Col: 13, Line: 3},
					},
				},
			},
			Else: Block{
				start: Pos{Col: 8, Line: 4},
				end:   Pos{Col: 2, Line: 6},
				Statements: []Statement{
					Instruction{
						Text:   "else this",
						quoted: `"else this"`,
						start:  Pos{Col: 2, Line: 5},
						end:    Pos{Col: 13, Line: 5},
					},
				},
			},
		},
	}})
	check.Eq(t, s.Statements[0].End(), Pos{Col: 2, Line: 6})
}

func TestSwitchCanBeEmpty(t *testing.T) {
	s, err := ParseString(`
switch "thing" {}
`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Switch{
			start: Pos{Col: 1, Line: 2},
			end:   Pos{Col: 18, Line: 2},
			Subject: String{
				Text:   "thing",
				quoted: `"thing"`,
				start:  Pos{Col: 8, Line: 2},
				end:    Pos{Col: 15, Line: 2},
			},
		},
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
			start: Pos{Col: 1, Line: 2},
			end:   Pos{Col: 2, Line: 5},
			Subject: String{
				Text:   "x",
				quoted: `"x"`,
				start:  Pos{Col: 8, Line: 2},
				end:    Pos{Col: 11, Line: 2},
			},
			Cases: []SwitchCase{
				{
					Condition: String{
						Text:   "1",
						quoted: `"1"`,
						start:  Pos{Col: 7, Line: 3},
						end:    Pos{Col: 10, Line: 3},
					},
					Block: Block{
						start: Pos{Col: 11, Line: 3},
						end:   Pos{Col: 13, Line: 3},
					},
				},
				{
					Condition: String{
						Text:   "2",
						quoted: `"2"`,
						start:  Pos{Col: 7, Line: 4},
						end:    Pos{Col: 10, Line: 4},
					},
					Block: Block{
						start: Pos{Col: 11, Line: 4},
						end:   Pos{Col: 20, Line: 4},
						Statements: []Statement{
							Instruction{
								Text:   "two",
								quoted: `"two"`,
								start:  Pos{Col: 13, Line: 4},
								end:    Pos{Col: 18, Line: 4},
							},
						},
					},
				},
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
			start: Pos{Col: 1, Line: 2},
			end:   Pos{Col: 2, Line: 5},
			Subject: String{
				Text:   "x",
				quoted: `"x"`,
				start:  Pos{Col: 8, Line: 2},
				end:    Pos{Col: 11, Line: 2},
			},
			Cases: []SwitchCase{
				{
					Condition: String{
						Text:   "1",
						quoted: `"1"`,
						start:  Pos{Col: 7, Line: 3},
						end:    Pos{Col: 10, Line: 3},
					},
					Block: Block{
						start: Pos{Col: 11, Line: 3},
						end:   Pos{Col: 13, Line: 3},
					},
				},
				{
					IsDefault: true,
					Block: Block{
						start: Pos{Col: 15, Line: 4},
						end:   Pos{Col: 17, Line: 4},
					},
				},
			},
		},
	}})
}

func TestInfiniteLoopHasNoCondition(t *testing.T) {
	s, err := ParseString(`while { "do" }`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		InfiniteLoop{
			start: Pos{Col: 1, Line: 1},
			Block: Block{
				start: Pos{Col: 7, Line: 1},
				end:   Pos{Col: 15, Line: 1},
				Statements: []Statement{
					Instruction{
						Text:   "do",
						quoted: `"do"`,
						start:  Pos{Col: 9, Line: 1},
						end:    Pos{Col: 13, Line: 1},
					},
				},
			},
		},
	}})
	check.Eq(t, s.Statements[0].End(), Pos{Col: 15, Line: 1})
}

func TestWhileLoopHasConditionNextToTheWhileKeyword(t *testing.T) {
	s, err := ParseString(`while "condition" { "do" }`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		While{
			start: Pos{Col: 1, Line: 1},
			Condition: String{
				Text:   "condition",
				quoted: `"condition"`,
				start:  Pos{Col: 7, Line: 1},
				end:    Pos{Col: 18, Line: 1},
			},
			Block: Block{
				start: Pos{Col: 19, Line: 1},
				end:   Pos{Col: 27, Line: 1},
				Statements: []Statement{
					Instruction{
						Text:   "do",
						quoted: `"do"`,
						start:  Pos{Col: 21, Line: 1},
						end:    Pos{Col: 25, Line: 1},
					},
				},
			},
		},
	}})
	check.Eq(t, s.Statements[0].End(), Pos{Col: 27, Line: 1})
}

func TestDoWhileLoopHasConditionInFooter(t *testing.T) {
	s, err := ParseString(`do { "do" } while "condition"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		DoWhile{
			start: Pos{Col: 1, Line: 1},
			Block: Block{
				start: Pos{Col: 4, Line: 1},
				end:   Pos{Col: 12, Line: 1},
				Statements: []Statement{
					Instruction{
						Text:   "do",
						quoted: `"do"`,
						start:  Pos{Col: 6, Line: 1},
						end:    Pos{Col: 10, Line: 1},
					},
				},
			},
			Condition: String{
				Text:   "condition",
				quoted: `"condition"`,
				start:  Pos{Col: 19, Line: 1},
				end:    Pos{Col: 30, Line: 1},
			},
		},
	}})
	check.Eq(t, s.Statements[0].End(), Pos{Col: 30, Line: 1})
}

func TestLoopsCanHaveBreaks(t *testing.T) {
	s, err := ParseString(`while { break "destination" }`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		InfiniteLoop{
			start: Pos{Col: 1, Line: 1},
			Block: Block{
				start: Pos{Col: 7, Line: 1},
				end:   Pos{Col: 30, Line: 1},
				Statements: []Statement{
					Break{
						Text:   "destination",
						quoted: `"destination"`,
						start:  Pos{Col: 9, Line: 1},
						end:    Pos{Col: 28, Line: 1},
					},
				},
			},
		},
	}})
	check.Eq(t, s.Statements[0].End(), Pos{Col: 30, Line: 1})
}

func TestCallBlockHasOneStringInstruction(t *testing.T) {
	s, err := ParseString(`call "instruction"`)
	check.Eq(t, err, nil)
	check.Eq(t, s, &Structogram{Statements: []Statement{
		Call{
			Text:   "instruction",
			quoted: `"instruction"`,
			start:  Pos{Col: 1, Line: 1},
			end:    Pos{Col: 19, Line: 1},
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
			start: Pos{Col: 1, Line: 2},
			end:   Pos{Col: 2, Line: 10},
			Blocks: []Block{
				Block{
					start: Pos{Col: 2, Line: 3},
					end:   Pos{Col: 3, Line: 5},
					Statements: []Statement{
						Instruction{
							Text:   "block 1",
							quoted: `"block 1"`,
							start:  Pos{Col: 3, Line: 4},
							end:    Pos{Col: 12, Line: 4},
						},
					},
				},
				Block{
					start: Pos{Col: 2, Line: 6},
					end:   Pos{Col: 4, Line: 6},
				},
				Block{
					start: Pos{Col: 2, Line: 7},
					end:   Pos{Col: 3, Line: 9},
					Statements: []Statement{
						Instruction{
							Text:   "block 3",
							quoted: `"block 3"`,
							start:  Pos{Col: 3, Line: 8},
							end:    Pos{Col: 12, Line: 8},
						},
					},
				},
			},
		},
		Parallel{
			start: Pos{Col: 1, Line: 12},
			end:   Pos{Col: 12, Line: 12},
		},
	}})
}

func TestIncompleteSwitchGivesParseError(t *testing.T) {
	_, err := ParseString(`switch "" {`)
	check.Eq(t, err.Error(), "parse error: token '}' expected")
}
