package parser

type Structogram struct {
	Title      String
	Statements []Statement
}

type Statement interface {
	Start() Pos
	End() Pos
}

type Pos struct {
	Col, Line int
}

type String struct {
	Text   string
	quoted string
	start  Pos
	end    Pos
}

func (s String) Start() Pos { return s.start }
func (s String) End() Pos   { return s.end }

type Instruction struct {
	Text   string
	quoted string
	start  Pos
	end    Pos
}

func (i Instruction) Start() Pos { return i.start }
func (i Instruction) End() Pos   { return i.end }

type Break struct {
	Text   string
	quoted string
	start  Pos
	end    Pos
}

func (b Break) Start() Pos { return b.start }
func (b Break) End() Pos   { return b.end }

type Call struct {
	Text   string
	quoted string
	start  Pos
	end    Pos
}

func (c Call) Start() Pos { return c.start }
func (c Call) End() Pos   { return c.end }

type Block struct {
	Statements []Statement
	start      Pos
	end        Pos
}

func (b Block) Start() Pos { return b.start }
func (b Block) End() Pos   { return b.end }

type If struct {
	Condition String
	Then      Block
	start     Pos
}

func (i If) Start() Pos { return i.start }
func (i If) End() Pos   { return i.Then.End() }

type IfElse struct {
	Condition String
	Then      Block
	Else      Block
	start     Pos
}

func (i IfElse) Start() Pos { return i.start }
func (i IfElse) End() Pos   { return i.Else.End() }

type Switch struct {
	Subject string
	Cases   []SwitchCase
}

func (Switch) Start() Pos { return Pos{} } // TODO
func (Switch) End() Pos   { return Pos{} } // TODO

type SwitchCase struct {
	// IsDefault and Condition are exclusive, a switch case is either the
	// default (case default {...}) or has a Condition (case "condition" {...}).
	IsDefault bool
	Condition string
	Block     Block
}

type Parallel struct {
	Blocks []Block
}

func (Parallel) Start() Pos { return Pos{} } // TODO
func (Parallel) End() Pos   { return Pos{} } // TODO

type InfiniteLoop struct {
	Block Block
}

func (InfiniteLoop) Start() Pos { return Pos{} } // TODO
func (InfiniteLoop) End() Pos   { return Pos{} } // TODO

type While struct {
	Condition string
	Block     Block
}

func (While) Start() Pos { return Pos{} } // TODO
func (While) End() Pos   { return Pos{} } // TODO

type DoWhile struct {
	Block     Block
	Condition string
}

func (DoWhile) Start() Pos { return Pos{} } // TODO
func (DoWhile) End() Pos   { return Pos{} } // TODO
