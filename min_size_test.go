package main

import (
	"testing"

	"github.com/gonutz/check"

	"github.com/gonutz/structorama/parser"
)

func TestInstructionHasTextSizePlusMargins(t *testing.T) {
	// 	 _______________
	// 	|               |
	// 	|  Instruction  |
	// 	|_______________|
	//
	// Instructions are just text in a rectangle, they have a margin around the
	// text which is half the line height on all sides. Empty instructions are
	// thus line height by line height in size.
	checkMinSize(t,
		&mockPainter{lineHeight: 10},
		parser.Instruction(""),
		5+5, 5+5,
	)
	checkMinSize(t,
		&mockPainter{lineHeight: 10, textW: 100, textH: 20},
		parser.Instruction("100x20 text"),
		5+100+5, 5+20+5,
	)
}

func checkMinSize(t *testing.T, p painter, x interface{}, wantW, wantH int) {
	t.Helper()
	w, h := minSize(p, x)
	check.Eq(t, w, wantW, "width")
	check.Eq(t, h, wantH, "height")
}

func TestCallHasVerticalBarsAroundText(t *testing.T) {
	// 	 __________
	// 	| |      | |
	// 	| | Call | |
	// 	|_|______|_|
	//
	// A call is like an instruction but it has additional vertical bars on the
	// left and right edges. The bars are inset by half the line height. The
	// text has a margin of halt the line height on all sides, just like
	// instructions. Thus a call is always at least line height high and 2*line
	// height wide.
	checkMinSize(t,
		&mockPainter{lineHeight: 10},
		parser.Call(""),
		5+1+10+1+5, 10,
	)
	checkMinSize(t,
		&mockPainter{lineHeight: 10, textW: 50, textH: 15},
		parser.Call("call"),
		5+1+5+50+5+1+5, 5+15+5,
	)
}

func TestBreakHasFourthOfHeightWideArrowOnTheLeft(t *testing.T) {
	// 	 ______________
	// 	| /            |
	// 	|/   Break     |
	// 	|\   Statement |
	// 	|_\____________|
	//
	// 	|_|
	// 	height/4
	//
	// Break boxes have an arrow pointing left at the left edge. It spans the
	// whole box's height and we want it to be height/4 wide at the top and
	// bottom. The text is surrounded by a margin of half the line height, like
	// instructions, the left margin starts where the arrow ends.
	checkMinSize(t,
		&mockPainter{lineHeight: 16},
		parser.Break(""),
		// Box height is 16=2*half line height, this means box width is 16/4=4
		// for the arrow plus the text margin, again 2*half line height.
		4+1+8+8, 8+8,
	)
	checkMinSize(t,
		&mockPainter{lineHeight: 16, textW: 40, textH: 20},
		parser.Break("break"),
		// Box height is 8+20+8=36 so the arrow is 36/4=9 wide.
		9+1+8+40+8, 8+20+8,
	)
}

func TestBlockWidthIsMaxOfPartsHeightIsSumOfParts(t *testing.T) {
	w, h := minSizeBlock(10, []size{{50, 20}, {100, 30}, {40, 40}})
	check.Eq(t, w, 100)
	check.Eq(t, h, 20+1+30+1+40)
}

func TestBlockIsAtLeastMarginInSize(t *testing.T) {
	w, h := minSizeBlock(10, nil)
	check.Eq(t, w, 10)
	check.Eq(t, h, 10)

	w, h = minSizeBlock(20, []size{{1, 1}})
	check.Eq(t, w, 20)
	check.Eq(t, h, 20)
}

func TestBlockHeightIsNeverNegative(t *testing.T) {
	w, h := minSizeBlock(0, nil)
	check.Eq(t, w, 0)
	check.Eq(t, h, 0)
}

// 	 _________________ 	 __________________ 	 __________________
// 	|  _______________|	| condition up top |	|  |               |
// 	| |               |	|   _______________|	|  | do-while loop |
// 	| | infinite loop |	|  |               |	|  |_______________|
// 	| |_______________|	|  | while loop    |	| bottom condition |
// 	|_________________|	|__|_______________|	|__________________|
//
// Loops are a rectangle block inside the outer rectangle with an optional
// condition either at the top or bottom. The inner rectangle is always line
// height right of the left edge. The top and bottom inserts are line height
// at minimum, or the text plus half the line height on all sides.

func TestEmptyInfiniteLoopIsAtLeastTheSizeOfMargins(t *testing.T) {
	w, h := minSizeInfiniteLoop(10, 0, 0)
	check.Eq(t, w, 10+1+10)
	check.Eq(t, h, 10+1+10+1+10)

	w, h = minSizeInfiniteLoop(20, 100, 50)
	check.Eq(t, w, 20+1+100)
	check.Eq(t, h, 20+1+50+1+20)
}

func TestLongWhileConditionDominatesWidth(t *testing.T) {
	// The minimum sizes for the parts are the margin.
	w, h := minSizeWhile(10, 0, 0, 0, 0)
	check.Eq(t, w, 10+1+10)
	check.Eq(t, h, 10+1+10)

	// Here the block size dominates the width.
	w, h = minSizeWhile(10, 40, 15, 100, 50)
	check.Eq(t, w, 10+1+100)
	check.Eq(t, h, 5+15+5+1+50)

	// Here the condition dominates the width.
	w, h = minSizeWhile(20, 200, 15, 100, 50)
	check.Eq(t, w, 200)
	check.Eq(t, h, 10+15+10+1+50)
}

func TestParallelIsSumOfBlocksWideAndMaxBlockHigh(t *testing.T) {
	// 	 ______________
	// 	| /          \ |
	// 	|/____________\|
	// 	|    |    |    |
	// 	|    |    |    |
	// 	|____|____|____|
	// 	|\            /|
	// 	|_\__________/_|
	//
	// Every part is at least the size of margin. There is a one pixel high line
	// between the top, center and bottom parts.
	w, h := minSizeParallel(10, nil)
	check.Eq(t, w, 10+10+10)
	check.Eq(t, h, 10+1+10+1+10)

	// Not only for empty Parallels, also for tiny blocks, the minimum size of
	// each part is margin.
	w, h = minSizeParallel(10, []size{{1, 1}})
	check.Eq(t, w, 10+10+10)
	check.Eq(t, h, 10+1+10+1+10)

	w, h = minSizeParallel(10, []size{{50, 20}})
	check.Eq(t, w, 50)
	check.Eq(t, h, 10+1+20+1+10)

	// Blocks have a one pixel wide vertical line between them.
	w, h = minSizeParallel(10, []size{{50, 20}, {100, 10}})
	check.Eq(t, w, 50+1+100)
	check.Eq(t, h, 10+1+20+1+10)

	w, h = minSizeParallel(10, []size{{50, 20}, {100, 10}, {20, 40}})
	check.Eq(t, w, 50+1+100+1+20)
	check.Eq(t, h, 10+1+40+1+10)
}
