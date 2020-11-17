package parser

import (
	"fmt"
	"unicode"
)

func tokenize(code string) ([]token, error) {
	var tokens []token
	var err error
	runes := []rune(code)

	// pos is the current index into runes, col and line are 1-indexed position
	// info for error messages.
	pos, col, line := 0, 1, 1

	// EOF indicates the end of file.
	const EOF = 0

	makeErr := func(msg string) error {
		return fmt.Errorf("%d:%d: %s", line, col, msg)
	}

	cur := func() rune {
		if pos < len(runes) {
			return runes[pos]
		}
		return EOF
	}
	next := func() {
		if pos < len(runes) {
			col++
			if runes[pos] == '\n' {
				col = 1
				line++
			}
			pos++
		}
	}

	// start... remember where the current token started, next() will update the
	// current position.
	startPos := 0
	startCol := col
	startLine := line
	emit := func(typ tokenType) {
		tokens = append(tokens, token{
			typ:  typ,
			text: string(runes[startPos:pos]),
			col:  startCol,
			line: startLine,
		})
		startPos = pos
		startCol = col
		startLine = line
	}

	// The main tokenize loop uses only the helper functions declared above.
	for cur() != EOF {
		switch cur() {
		case '{', '}':
			typ := tokenType(cur())
			next()
			emit(typ)
		case '"':
			next()
			for {
				if cur() == '"' {
					next() // Skip the closing quote.
					break
				} else if cur() == '\\' {
					// Escape sequence "\\", "\"" or "\n".
					next()
					if cur() == '\\' || cur() == 'n' || cur() == '"' {
						next()
					} else if cur() == EOF {
						return nil, makeErr("unexpected end of input after '\\' in string literal")
					} else {
						return nil, makeErr("unknown escape sequence (only 'n', '\\' and '\"' can follow after '\\')")
					}
				} else if cur() == EOF {
					return nil, makeErr("unexpected end of input in string literal")
				} else {
					next()
				}
			}
			emit(tokenString)
		default:
			if unicode.IsSpace(cur()) {
				for unicode.IsSpace(cur()) {
					next()
				}
				emit(tokenSpace)
			} else if unicode.IsLetter(cur()) {
				for unicode.IsLetter(cur()) {
					next()
				}
				emit(tokenID)
			} else {
				return nil, makeErr("illegal character " + fmt.Sprintf("%q", cur()))
			}
		}
	}

	return tokens, err
}

type token struct {
	typ  tokenType
	text string
	col  int
	line int
}

type tokenType int

const (
	tokenID = iota
	tokenString
	tokenSpace
)

func (t tokenType) String() string {
	switch t {
	case tokenID:
		return "identifier"
	case tokenString:
		return "string"
	case tokenSpace:
		return "white space"
	default:
		return fmt.Sprintf("token %q", rune(t))
	}
}
