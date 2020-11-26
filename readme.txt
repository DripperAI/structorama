TODO
----

create a language that can handle all cases from https://de.wikipedia.org/wiki/Nassi-Shneiderman-Diagramm ✓ (see below)
build a parser for the language that outputs an AST ✓
build a formatter for the AST ✓
build a GUI with a text editor and a graphical preview ✓
	have load/save for files
	start the program with a default showcase diagram containing all constructs ✓
	have a help, some way to know the syntax of the language ✓
	generate PDFs ✓
have a PDF generator with input AST and output .pdf file ✓


Syntax
------

Example:
--------------------------------------------
title "optional diagram caption"

"counter := 0"

if "only if" {
}

if "if-else" "T" {
} else "F" {
}

switch "subject" {
	case "1" {}
	case "2" {}
	case default {}
}

while {
	"infinite loop"
}

while "i=0; i<10; i++" {
	break "early exit the loop"
}

do {} while "i<10"

call "some function"

parallel {
	{
		if "nested things" {
			"in block 1"
		}
	}
	{}
	{
		"block right of the empty block"
	}
}
--------------------------------------------

More realistic example:
--------------------------------------------
title "counter example (countering what?)"

"counter = 0"
while "counter != 10" {
	"print \"counter is\""

	if "counter % 2 == 0" "yes" {
		call "printEven()"
	} else "no" {
		call "printOdd()"
	}

	"counter++"
}
--------------------------------------------
