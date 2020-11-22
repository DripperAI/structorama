package parser

import "testing"

func TestFormatter(t *testing.T) {
	code, err := FormatString(`
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
	`)
	if err != nil {
		t.Fatal(err)
	}
	want := `title "some caption"

"instrucion"
"with \n and tabs:
	and real line breaks"

call "multiple empty lines turn into only two"

break "while single empty lines remain"
"and no empty lines are fine too"
"but at least one line break is printed"
break "can have multiple lines"
"which is OK"
`
	if code != want {
		t.Errorf("have\n---\n%s\n---\nbut want\n---\n%s\n---", code, want)
	}
}

func TestThereAreAlwaysOneEmptyLineBelowTheTitle(t *testing.T) {
	code, err := FormatString(`title "some caption" "instruction"`)
	if err != nil {
		t.Fatal(err)
	}
	want := `title "some caption"

"instruction"
`
	if code != want {
		t.Errorf("have\n---\n%s\n---\nbut want\n---\n%s\n---", code, want)
	}
}

func TestThereDoesNotHaveToBeATitle(t *testing.T) {
	code, err := FormatString(`
	
	
	
	"instruction"`)
	if err != nil {
		t.Fatal(err)
	}
	want := `"instruction"
`
	if code != want {
		t.Errorf("have\n---\n%s\n---\nbut want\n---\n%s\n---", code, want)
	}
}

func TestOneEmptyLineIsKeptBetweenStatements(t *testing.T) {
	code, err := FormatString(`"a"

"b"
`)
	if err != nil {
		t.Fatal(err)
	}
	want := `"a"

"b"
`
	if code != want {
		t.Errorf("have\n---\n%s\n---\nbut want\n---\n%s\n---", code, want)
	}
}

func TestIfElseIsFormattedLikeGo(t *testing.T) {
	code, err := FormatString(`if"true"{"then"}else{"not"}`)
	if err != nil {
		t.Fatal(err)
	}
	want := `if "true" {
	"then"
} else {
	"not"
}
`
	if code != want {
		t.Errorf("have\n---\n%s\n---\nbut want\n---\n%s\n---", code, want)
	}
}
