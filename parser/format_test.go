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

func checkFormatting(t *testing.T, original, want string) {
	have, err := FormatString(original)
	if err != nil {
		t.Fatal(err)
	}
	if have != want {
		t.Errorf("have\n---\n%s\n---\nbut want\n---\n%s\n---", have, want)
	}
}
