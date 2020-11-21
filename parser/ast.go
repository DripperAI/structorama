package parser

type Structogram struct {
	Title      string
	Statements []Statement
}

type Statement interface{}

type Instruction string

type Break string

type Call string

type Block []Statement

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
