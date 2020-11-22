package main

import (
	"fmt"
	"strings"

	"github.com/gonutz/wui"

	"github.com/gonutz/structorama/parser"
)

// TODO Have a setting for the text for if's true and false cases and for a
// switch's default case.

func main() {
	codeFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Courier New",
		Height: -19,
		Bold:   true,
	})
	previewFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -17,
	})

	window := wui.NewWindow()
	window.SetFont(codeFont)
	window.SetTitle("Structorama")

	codeEditor := wui.NewTextEdit()
	codeEditor.SetAnchors(wui.AnchorMinAndCenter, wui.AnchorMinAndMax)
	codeEditor.SetSize(300, 400)
	codeEditor.SetWritesTabs(true)
	window.Add(codeEditor)

	preview := wui.NewPaintBox()
	// TODO
	//preview.SetFont(previewFont)
	preview.SetAnchors(wui.AnchorMaxAndCenter, wui.AnchorMinAndMax)
	preview.SetX(300)
	preview.SetSize(300, 400)
	window.Add(preview)

	window.SetState(wui.WindowMaximized)
	window.SetOnShow(codeEditor.Focus)

	const example = `title "optional diagram caption"

"counter := 0"

if "only if" {
}

if "if-else" {
} else {
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
}`

	// TODO
	//codeEdit.SetLineBreak("\n"), this should probably be the default in Go.
	codeEditor.SetText(strings.Replace(example, "\n", "\r\n", -1))

	var lastValidStructogram *parser.Structogram
	preview.SetOnPaint(func(canvas *wui.Canvas) {
		canvas.SetFont(previewFont)
		canvas.FillRect(
			0, 0, canvas.Width(), canvas.Height(),
			wui.RGB(255, 255, 255),
		)

		s, err := parser.ParseString(codeEditor.Text())
		if err == nil {
			lastValidStructogram = s
		}

		paintStructogram(
			offsetPainter{p: canvasPainter{c: canvas}, dx: 10, dy: 10},
			lastValidStructogram,
		)

		if err != nil {
			canvas.TextRectFormat(
				0, 0, canvas.Width(), canvas.Height(),
				err.Error(), wui.FormatCenter, wui.RGB(255, 0, 0),
			)
		}
	})

	formatCode := func() {
		code, err := parser.FormatString(codeEditor.Text())
		if err == nil {
			codeEditor.SetText(strings.Replace(code, "\n", "\r\n", -1))
		} else {
			wui.MessageBoxError("Formatting Error", err.Error())
		}
	}

	codeEditor.SetOnTextChange(preview.Paint)

	window.SetShortcut(formatCode, wui.KeyControl, wui.KeyF)
	window.SetShortcut(window.Close, wui.KeyEscape)

	window.Show()
}

type painter interface {
	// Text paints text s aligned with its top-left at (x,y).
	Text(x, y int, s string)
	// TestSize returns the size of the given text without any margins.
	TextSize(s string) (width, height int)
	// Rect paints a one pixel wide rectangle of the given width and height. The
	// width and height include the outer pixels that are painted.
	Rect(x, y, width, height int)
	// Line paints a one pixel wide line from (x1,y1) to (x2,y2), including both
	// end points.
	Line(x1, y1, x2, y2 int)
	// LineHeight is the height of one line of text.
	LineHeight() int
}

type canvasPainter struct {
	c *wui.Canvas
}

const infinite = 0x0FFFFFFF

func (p canvasPainter) Text(x, y int, s string) {
	p.c.TextRect(x, y, infinite, infinite, s, wui.RGB(0, 0, 0))
}

func (p canvasPainter) TextSize(s string) (width, height int) {
	return p.c.TextRectExtent(s, infinite)
}

func (p canvasPainter) Rect(x, y, width, height int) {
	p.c.DrawRect(x, y, width, height, wui.RGB(0, 0, 0))
}

func (p canvasPainter) Line(x1, y1, x2, y2 int) {
	p.c.Line(x1, y1, x2, y2, wui.RGB(0, 0, 0))
	// Draw the last pixel, Canvas.Line does not include it.
	p.c.Line(x2, y2, x2+1, y2, wui.RGB(0, 0, 0))
}

func (p canvasPainter) LineHeight() int {
	_, h := p.c.TextExtent("|")
	return h
}

type offsetPainter struct {
	p  painter
	dx int
	dy int
}

func (p offsetPainter) Text(x, y int, s string) {
	p.p.Text(x+p.dx, y+p.dy, s)
}

func (p offsetPainter) TextSize(s string) (width, height int) {
	return p.p.TextSize(s)
}

func (p offsetPainter) Rect(x, y, width, height int) {
	p.p.Rect(x+p.dx, y+p.dy, width, height)
}

func (p offsetPainter) Line(x1, y1, x2, y2 int) {
	p.p.Line(x1+p.dx, y1+p.dy, x2+p.dx, y2+p.dy)
}

func (p offsetPainter) LineHeight() int {
	return p.p.LineHeight()
}

func paintStructogram(p painter, x *parser.Structogram) {
	if x.Title.Text != "" {
		p.Text(0, 0, x.Title.Text)
		_, h := p.TextSize(x.Title.Text)
		p = offsetPainter{p: p, dy: h + 5}
	}
	body := parser.Block{Statements: x.Statements}
	width, height := minSize(p, body)
	p.Rect(-1, -1, width+2, height+2)
	paintIn(p, body, width, height)
}

func paintIn(p painter, node interface{}, width, height int) {
	margin := p.LineHeight()
	switch x := node.(type) {

	case parser.Instruction:
		p.Text(margin/2, margin/2, x.Text)

	case parser.Call:
		left := margin / 2
		right := width - 1 - left
		p.Line(left, 0, left, height-1)
		p.Line(right, 0, right, height-1)
		p.Text(left+1+margin/2, margin/2, x.Text)

	case parser.Break:
		p.Line(0, (height-1)/2, height/4, 0)
		p.Line(0, height/2, height/4, height-1)
		p.Text(height/4+1+margin/2, margin/2, x.Text)

	case parser.Block:
		sizes := make([]size, len(x.Statements))
		for i := range sizes {
			sizes[i].width, sizes[i].height = minSize(p, x.Statements[i])
		}
		areas := blockPaintAreas(width, height, sizes)
		for i := range x.Statements {
			if i > 0 {
				y := areas[i].y - 1
				p.Line(0, y, width-1, y)
			}
			paintIn(
				offsetPainter{p: p, dx: areas[i].x, dy: areas[i].y},
				x.Statements[i],
				areas[i].width,
				areas[i].height,
			)
		}

	case parser.InfiniteLoop:
		area := paintInfiniteLoopLines(p, x, width, height)
		paintIn(
			offsetPainter{p: p, dx: area.x, dy: area.y},
			x.Block,
			area.width,
			area.height,
		)

	case parser.While:
		area := paintWhileLoop(p, x, width, height)
		paintIn(
			offsetPainter{p: p, dx: area.x, dy: area.y},
			x.Block,
			area.width,
			area.height,
		)

	case parser.DoWhile:
		area := paintDoWhileLoop(p, x, width, height)
		paintIn(
			offsetPainter{p: p, dx: area.x, dy: area.y},
			x.Block,
			area.width,
			area.height,
		)

	case parser.Parallel:
		sizes := make([]size, len(x.Blocks))
		for i := range sizes {
			sizes[i].width, sizes[i].height = minSize(p, x.Blocks[i])
		}
		areas := paintParallel(p, sizes, width, height)
		for i, block := range x.Blocks {
			paintIn(
				offsetPainter{p: p, dx: areas[i].x, dy: areas[i].y},
				block,
				areas[i].width,
				areas[i].height,
			)
		}

	case parser.If:
		// An If is the same as an IfElse with an empty Else.
		paintIn(p, parser.IfElse{
			Condition: parser.String{Text: x.Condition},
			Then:      x.Then,
		}, width, height)

	case parser.IfElse:
		// TODO

	case parser.Switch:
		// TODO

	default:
		panic("TODO paintIn: unhandled structogram node: " +
			fmt.Sprintf("%T", node))
	}
}

// minSize returns the minimum inner size of the given node when painted with
// the given painter. No border around the node is considered. Parent nodes are
// expeced to take them into account. See the accompanying unit tests for ASCII
// art and explanation of these sizes.
func minSize(p painter, node interface{}) (width, height int) {
	margin := p.LineHeight()
	switch x := node.(type) {

	case parser.Instruction:
		textW, textH := p.TextSize(x.Text)
		return textW + margin, textH + margin

	case parser.Call:
		textW, textH := p.TextSize(x.Text)
		return 2*margin + 2 + textW, textH + margin

	case parser.Break:
		textW, textH := p.TextSize(x.Text)
		height := margin + textH
		width := height/4 + 1 + textW + margin
		return width, height

	case parser.Block:
		sizes := make([]size, len(x.Statements))
		for i := range sizes {
			sizes[i].width, sizes[i].height = minSize(p, x.Statements[i])
		}
		return minSizeBlock(margin, sizes)

	case parser.InfiniteLoop:
		blockW, blockH := minSize(p, x.Block)
		return minSizeInfiniteLoop(margin, blockW, blockH)

	case parser.While:
		textW, textH := p.TextSize(x.Condition)
		blockW, blockH := minSize(p, x.Block)
		return minSizeWhile(margin, textW, textH, blockW, blockH)

	case parser.DoWhile:
		// While and DoWhile have the same size, one has the block up top, the
		// other one at the bottom.
		return minSize(p, parser.While{
			Condition: x.Condition,
			Block:     x.Block,
		})

	case parser.Parallel:
		sizes := make([]size, len(x.Blocks))
		for i := range sizes {
			sizes[i].width, sizes[i].height = minSize(p, x.Blocks[i])
		}
		return minSizeParallel(margin, sizes)

	case parser.If:
		// An If is the same as an IfElse with an empty Else.
		return minSize(p, parser.IfElse{
			Condition: parser.String{Text: x.Condition},
			Then:      x.Then,
		})

	case parser.IfElse:
		return 100, 30 // TODO

	case parser.Switch:
		return 100, 30 // TODO

	default:
		panic("TODO minSize: unhandled structogram node: " +
			fmt.Sprintf("%T", node))
	}
}

type size struct {
	width, height int
}

func minSizeBlock(margin int, sizes []size) (width, height int) {
	var maxWidth, totalHeight int
	for _, s := range sizes {
		if maxWidth < s.width {
			maxWidth = s.width
		}
		totalHeight += s.height
	}
	return max(margin, maxWidth), max(margin, totalHeight+len(sizes)-1)
}

func minSizeParallel(margin int, sizes []size) (width, height int) {
	for _, s := range sizes {
		width += s.width
		if height < s.height {
			height = s.height
		}
	}
	width += len(sizes) - 1
	height += 2*margin + 2
	return max(3*margin, width), max(3*margin+2, height)
}

func minSizeInfiniteLoop(margin, blockW, blockH int) (width, height int) {
	width = margin + 1 + max(margin, blockW)
	height = 2*margin + 2 + max(margin, blockH)
	return
}

func minSizeWhile(margin, conditionW, conditionH, blockW, blockH int) (width, height int) {
	width = max(conditionW, margin+1+max(margin, blockW))
	height = margin + conditionH + 1 + max(margin, blockH)
	return
}

func blockPaintAreas(width, height int, blockSizes []size) []rectangle {
	r := make([]rectangle, len(blockSizes))
	y := 0
	for i := range r {
		r[i].x = 0
		r[i].y = y
		r[i].width = width
		r[i].height = blockSizes[i].height
		y += r[i].height + 1
	}
	if len(r) == 1 {
		r[len(r)-1].height = height
	}
	if len(r) >= 2 {
		r[len(r)-1].height = height - (r[len(r)-2].y + r[len(r)-2].height + 1)
	}
	return r
}

func paintInfiniteLoopLines(p painter, loop parser.InfiniteLoop, width, height int) (blockArea rectangle) {
	margin := p.LineHeight()
	bottom := height - 1 - margin
	p.Line(margin, margin, width-1, margin)
	p.Line(margin, margin, margin, bottom)
	p.Line(margin, bottom, width-1, bottom)
	return rectangle{
		x:      margin + 1,
		y:      margin + 1,
		width:  width - margin - 1,
		height: height - 2 - 2*margin,
	}
}

func paintWhileLoop(p painter, while parser.While, width, height int) (blockArea rectangle) {
	margin := p.LineHeight()
	_, textH := p.TextSize(while.Condition)
	top := margin + textH
	p.Text(margin/2, margin/2, while.Condition)
	p.Line(margin, top, width-1, top)
	p.Line(margin, top, margin, height-1)
	return rectangle{
		x:      margin + 1,
		y:      top + 1,
		width:  width - margin - 1,
		height: height - (top + 1),
	}
}

func paintDoWhileLoop(p painter, do parser.DoWhile, width, height int) (blockArea rectangle) {
	margin := p.LineHeight()
	_, textH := p.TextSize(do.Condition)
	bottom := height - 1 - margin - textH
	p.Line(margin, 0, margin, bottom)
	p.Line(margin, bottom, width-1, bottom)
	p.Text(margin/2, bottom+1+margin/2, do.Condition)
	return rectangle{
		x:      margin + 1,
		y:      0,
		width:  width - margin - 1,
		height: bottom,
	}
}

func paintParallel(p painter, blockSizes []size, width, height int) []rectangle {
	margin := p.LineHeight()
	p.Line(0, margin, width-1, margin)
	p.Line(0, height-1-margin, width-1, height-1-margin)
	p.Line(0, margin-1, margin-1, 0)
	p.Line(width-margin, 0, width-1, margin-1)
	p.Line(0, height-margin, margin-1, height-1)
	p.Line(width-margin, height-1, width-1, height-margin)

	areas := make([]rectangle, len(blockSizes))
	totalBlockW := 0
	for _, s := range blockSizes {
		totalBlockW += s.width
	}
	scale := 0.0
	if totalBlockW > 0 {
		scale = float64(width-(len(blockSizes)-1)) / float64(totalBlockW)
	}
	x := 0
	for i := range areas {
		if i > 0 {
			p.Line(x-1, margin+1, x-1, height-margin-1)
		}
		areas[i].x = x
		areas[i].y = margin + 1
		areas[i].width = int(float64(blockSizes[i].width)*scale + 0.5)
		areas[i].height = height - 2*margin - 2
		x += 1 + areas[i].width
	}
	return areas
}

type rectangle struct {
	x, y, width, height int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
