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
	TrueText  String
	Then      Block
	start     Pos
}

func (i If) Start() Pos { return i.start }
func (i If) End() Pos   { return i.Then.End() }

type IfElse struct {
	Condition String
	TrueText  String
	Then      Block
	FalseText String
	Else      Block
	start     Pos
}

func (i IfElse) Start() Pos { return i.start }
func (i IfElse) End() Pos   { return i.Else.End() }

type Switch struct {
	Subject String
	Cases   []SwitchCase
	start   Pos
	end     Pos
}

func (s Switch) Start() Pos { return s.start }
func (s Switch) End() Pos   { return s.end }

type SwitchCase struct {
	// IsDefault and Condition are exclusive, a switch case is either the
	// default (case default {...}) or has a Condition (case "condition" {...}).
	IsDefault bool
	Condition String
	Block     Block
}

type Parallel struct {
	Blocks []Block
	start  Pos
	end    Pos
}

func (p Parallel) Start() Pos { return p.start }
func (p Parallel) End() Pos   { return p.end }

type InfiniteLoop struct {
	Block Block
	start Pos
}

func (i InfiniteLoop) Start() Pos { return i.start }
func (i InfiniteLoop) End() Pos   { return i.Block.End() }

type While struct {
	Condition String
	Block     Block
	start     Pos
}

func (w While) Start() Pos { return w.start }
func (w While) End() Pos   { return w.Block.End() }

type DoWhile struct {
	Block     Block
	Condition String
	start     Pos
}

func (d DoWhile) Start() Pos { return d.start }
func (d DoWhile) End() Pos   { return d.Condition.End() }
