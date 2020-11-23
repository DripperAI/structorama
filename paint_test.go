package main

import (
	"testing"

	"github.com/gonutz/check"

	"github.com/gonutz/structorama/parser"
)

func TestPaintInstruction(t *testing.T) {
	// 	 _______________
	// 	|               |
	// 	|  Instruction  |
	// 	|_______________|
	p := &mockPainter{lineHeight: 10}
	paintIn(p, parser.Instruction{Text: "instruction"}, 0, 0)
	p.checkPainting(t, `Text(5, 5, "instruction")`)
}

func TestPaintCall(t *testing.T) {
	// 	 __________
	// 	| |      | |
	// 	| | Call | |
	// 	|_|______|_|
	p := &mockPainter{lineHeight: 10}
	paintIn(p, parser.Call{Text: "call"}, 100, 50)
	p.checkPainting(t,
		`Line(5, 0, 5, 49)`,
		`Line(94, 0, 94, 49)`,
		`Text(11, 5, "call")`,
	)
}

func TestPaintBreak(t *testing.T) {
	// 	 ______________
	// 	| /            |
	// 	|/   Break     |
	// 	|\             |
	// 	|_\____________|
	//
	// 	|_|
	// 	line height / 4
	//
	// Example: 10x4 area
	// 	.x........
	// 	x.........
	// 	x.........
	// 	.x........
	//
	// Example: 10x8 area
	// 	..x.......
	// 	.x........
	// 	.x........
	// 	x.........
	// 	x.........
	// 	.x........
	// 	.x........
	// 	..x.......
	p := &mockPainter{lineHeight: 10}
	paintIn(p, parser.Break{Text: "break"}, 100, 40)
	p.checkPainting(t,
		`Line(0, 19, 10, 0)`,
		`Line(0, 20, 10, 39)`,
		`Text(16, 5, "break")`,
	)

	p = &mockPainter{lineHeight: 10}
	paintIn(p, parser.Break{Text: "break"}, 100, 41)
	p.checkPainting(t,
		`Line(0, 20, 10, 0)`,
		`Line(0, 20, 10, 40)`,
		`Text(16, 5, "break")`,
	)
}

func TestPaintingEmptyBlockDoesNothing(t *testing.T) {
	p := &mockPainter{lineHeight: 10}
	paintIn(p, parser.Block{}, 100, 40)
	p.checkPainting(t)
}

func TestBlockPaintsPartsOnePixelApartOverWholeAvailableWidth(t *testing.T) {
	// If there is only one block, it always occupies the whole area.
	areas := blockPaintAreas(100, 200, []size{{50, 30}})
	check.Eq(t, areas, []rectangle{{0, 0, 100, 200}})

	// All blocks but the last occupy the full width and their minimal height.
	// The last block occupies the remaining space at the bottom. Below each
	// block is a separating one pixel high line.
	areas = blockPaintAreas(100, 200, []size{
		{50, 30},
		{80, 20},
	})
	check.Eq(t, areas, []rectangle{
		{0, 0, 100, 30},
		{0, 31, 100, 169},
	})
}

func TestPaintingInfiniteLoop(t *testing.T) {
	// 	 _________________
	// 	|  _______________|
	// 	| |               |
	// 	| | infinite loop |
	// 	| |_______________|
	// 	|_________________|
	p := &mockPainter{lineHeight: 10}
	area := paintInfiniteLoopLines(p, parser.InfiniteLoop{}, 200, 100)
	p.checkPainting(t,
		`Line(10, 10, 199, 10)`,
		`Line(10, 89, 199, 89)`,
		`Line(10, 10, 10, 89)`,
	)
	check.Eq(t, area, rectangle{11, 11, 189, 78})
}

func TestPaintingWhileLoop(t *testing.T) {
	// 	 __________________
	// 	| condition up top |
	// 	|   _______________|
	// 	|  |               |
	// 	|  | while loop    |
	// 	|__|_______________|
	p := &mockPainter{lineHeight: 10, textW: 50, textH: 20}
	area := paintWhileLoop(p, parser.While{Condition: parser.String{
		Text: "condition",
	}}, 200, 100)
	p.checkPainting(t,
		`Text(5, 5, "condition")`,
		`Line(10, 30, 199, 30)`,
		`Line(10, 30, 10, 99)`,
	)
	check.Eq(t, area, rectangle{11, 31, 189, 69})
}

func TestPaintingDoWhileLoop(t *testing.T) {
	// 	 __________________
	// 	|  |               |
	// 	|  | do-while loop |
	// 	|  |_______________|
	// 	| bottom condition |
	// 	|__________________|
	p := &mockPainter{lineHeight: 10, textW: 50, textH: 20}
	area := paintDoWhileLoop(p, parser.DoWhile{Condition: parser.String{
		Text: "condition",
	}}, 200, 100)
	p.checkPainting(t,
		`Line(10, 0, 10, 69)`,
		`Line(10, 69, 199, 69)`,
		`Text(5, 75, "condition")`,
	)
	check.Eq(t, area, rectangle{11, 0, 189, 69})
}

func TestParallelPaintsLinesBetweenBlocks(t *testing.T) {
	p := &mockPainter{lineHeight: 10}
	areas := paintParallel(p, nil, 200, 100)
	p.checkPainting(t,
		`Line(0, 9, 9, 0)`,
		`Line(190, 0, 199, 9)`,
		`Line(0, 90, 9, 99)`,
		`Line(190, 99, 199, 90)`,
		`Line(0, 10, 199, 10)`,
		`Line(0, 89, 199, 89)`,
	)
	check.Eq(t, areas, nil)

	// If there is only one block, it fills the whole space.
	p = &mockPainter{lineHeight: 10}
	areas = paintParallel(p, []size{{1, 1}}, 200, 100)
	p.checkPainting(t,
		`Line(0, 9, 9, 0)`,
		`Line(190, 0, 199, 9)`,
		`Line(0, 90, 9, 99)`,
		`Line(190, 99, 199, 90)`,
		`Line(0, 10, 199, 10)`,
		`Line(0, 89, 199, 89)`,
	)
	check.Eq(t, areas, []rectangle{{0, 11, 200, 78}})

	// If there are two blocks, the space is split horizontally so that each
	// gets its relative percentage.
	p = &mockPainter{lineHeight: 10}
	// Make it 201 for the extra pixel between the two blocks.
	// The left block gets 10/40 of the space.
	// The right block gets 30/40 of the space.
	areas = paintParallel(p, []size{{10, 10}, {30, 10}}, 201, 100)
	p.checkPainting(t,
		`Line(0, 9, 9, 0)`,
		`Line(191, 0, 200, 9)`,
		`Line(0, 90, 9, 99)`,
		`Line(191, 99, 200, 90)`,
		`Line(0, 10, 200, 10)`,
		`Line(0, 89, 200, 89)`,
		`Line(50, 11, 50, 89)`,
	)
	check.Eq(t, areas, []rectangle{
		{0, 11, 50, 78},
		{51, 11, 150, 78},
	})
}
