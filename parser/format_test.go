package parser

import "testing"

func TestFormatter(t *testing.T) {
	checkFormatting(t, `
 title  "some caption"

	"instrucion"
	  "with \n and tabs:
	and real line breaks" 	
 
 
 
 call "multiple empty lines turn into only two" 
 


  break "while single empty lines remain"  
"and no empty lines are fine too"  "but at least one line break is printed"
break


"can have multiple lines"
"which is OK"
	`,

		`title "some caption"

"instrucion"
"with \n and tabs:
	and real line breaks"

call "multiple empty lines turn into only two"

break "while single empty lines remain"
"and no empty lines are fine too"
"but at least one line break is printed"
break "can have multiple lines"
"which is OK"
`)
}

func TestEmptyTitleStays(t *testing.T) {
	checkFormatting(t, `title ""`, `title ""`)
}

func TestThereAreAlwaysOneEmptyLineBelowTheTitle(t *testing.T) {
	checkFormatting(t,
		`title "some caption" "instruction"`,

		`title "some caption"

"instruction"
`)
}

func TestThereDoesNotHaveToBeATitle(t *testing.T) {
	checkFormatting(t, `
	
	
	
	"instruction"`,

		`"instruction"
`)
}

func TestOneEmptyLineIsKeptBetweenStatements(t *testing.T) {
	checkFormatting(t, `"a"

"b"
`,

		`"a"

"b"
`)
}

func TestIfElseIsFormattedLikeGo(t *testing.T) {
	checkFormatting(t,
		`if"true"{"then"}else{"not"}`,

		`if "true" {
	"then"
} else {
	"not"
}
`)
}

func TestIfElseTrueAndFalseTextsAppearNextToTheKeywords(t *testing.T) {
	checkFormatting(t, `
if"true""T"{}
if"true""T"{}else"F"{}
`,

		`if "true" "T" {
	
}
if "true" "T" {
	
} else "F" {
	
}
`)
}

func TestEmptyIfElseTrueAndFalseTextsStay(t *testing.T) {
	checkFormatting(t, `
if"true"""{}
if"true"""{}else""{}
`,

		`if "true" "" {
	
}
if "true" "" {
	
} else "" {
	
}
`)
}

func TestIfElseCanBeNested(t *testing.T) {
	checkFormatting(t,
		`if"true"{if"false"{}else{}}else{"not"}`,

		`if "true" {
	if "false" {
		
	} else {
		
	}
} else {
	"not"
}
`)
}

func TestFormattedIfElseMayHaveEmptyLinesAroundIt(t *testing.T) {
	checkFormatting(t,
		`if""{}else{}

if""{}else{}

if""{}else{}`,

		`if "" {
	
} else {
	
}

if "" {
	
} else {
	
}

if "" {
	
} else {
	
}
`)
}
func TestIfIsFormattedLikeGo(t *testing.T) {
	checkFormatting(t,
		`if"a"{"do a"}  if"b"{"do b"}



if"c"{"do c"}`,

		`if "a" {
	"do a"
}
if "b" {
	"do b"
}

if "c" {
	"do c"
}
`)
}

func TestFormatSwitch(t *testing.T) {
	checkFormatting(t,
		`switch"what"{case"1"{"do 1"}case"2"{}case default{"else"}}


switch""{}`,

		`switch "what" {
	case "1" {
		"do 1"
	}
	case "2" {
		
	}
	case default {
		"else"
	}
}

switch "" {
	
}
`)
}

func TestFormatInfiniteLoop(t *testing.T) {
	checkFormatting(t,
		`while{"loop"}
		
	
	while{break "bbb"}`,

		`while {
	"loop"
}

while {
	break "bbb"
}
`)
}

func TestFormatWhileLoop(t *testing.T) {
	checkFormatting(t,
		`while   "a==b"{"loop"}
		
	
	while"1!=2"{break "bbb"}`,

		`while "a==b" {
	"loop"
}

while "1!=2" {
	break "bbb"
}
`)
}

func TestFormatDoWhileLoop(t *testing.T) {
	checkFormatting(t,
		`do{"loop"}while"false"
		
	
	do{"loop"}while"false"`,

		`do {
	"loop"
} while "false"

do {
	"loop"
} while "false"
`)
}

func TestFormatParallelBlocks(t *testing.T) {
	checkFormatting(t,
		`parallel{}parallel{{}}
	
	
	parallel{{"123"}{"456"}}`,

		`parallel {
	
}
parallel {
	{
		
	}
}

parallel {
	{
		"123"
	}
	{
		"456"
	}
}
`)
}

func TestBlockContentHasEmptyLines(t *testing.T) {
	checkFormatting(t, `
if "" {
	
	"a"
	"b"
	
	
	"c"
}`,

		`if "" {
	"a"
	"b"

	"c"
}
`)
}

func checkFormatting(t *testing.T, original, want string) {
	have, err := FormatString(original)
	if err != nil {
		t.Fatal(err)
	}
	if have != want {
		t.Errorf("have\n---\n%s\n---\nbut want\n---\n%s\n---", have, want)
	}
}
