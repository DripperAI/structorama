package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os/exec"
	"strings"

	"github.com/jung-kurt/gofpdf"

	"github.com/gonutz/gofont"
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

	const example2 = `title "optional diagram caption"

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

	const example = `if "ashudihasudishaiu" {
	call "asdi"
} else {
	call "asdi"
}
`

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
				err.Error(), wui.FormatBottomCenter, wui.RGB(255, 0, 0),
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

	exportPDF := func() {
		// Unfortunately implementing a pdfPainter using the gofpdf library
		// proved to be difficult. Instead we now just create a pixel-based
		// image, draw to it and render that into the PDF instead.

		fontPath := "C:/Windows/Fonts/" + previewFont.Desc.Name + ".ttf"
		font, err := gofont.LoadFromFile(fontPath)
		if err != nil {
			wui.MessageBoxError("Cannot load font", err.Error())
			return
		}
		font.HeightInPixels = 20

		// DIN A4 pages are 210 x 297 mm in size,we keep our image at the same
		// aspect ratio.
		img := image.NewRGBA(image.Rect(0, 0, 3*210, 3*297))
		paintStructogram(
			offsetPainter{p: imagePainter{img: img, font: font}, dx: 10, dy: 10},
			lastValidStructogram,
		)

		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		var buf bytes.Buffer
		png.Encode(&buf, img)
		pdf.RegisterImageOptionsReader(
			"diagram.png",
			gofpdf.ImageOptions{ImageType: "PNG"},
			bytes.NewReader(buf.Bytes()),
		)
		pdf.Image("diagram.png", 0, 0, 0, 0, false, "", 0, "")

		dlg := wui.NewFileSaveDialog()
		dlg.SetTitle("Select output path")
		dlg.AddFilter("PDF File", ".pdf")
		dlg.SetInitialPath("diagram.pdf")
		if ok, path := dlg.Execute(window); ok {
			err := pdf.OutputFileAndClose(path)
			if err != nil {
				wui.MessageBoxError("Error exporting PDF", err.Error())
			} else {
				exec.Command("cmd", "/C", path).Start()
			}
		}
	}

	codeEditor.SetOnTextChange(preview.Paint)

	window.SetShortcut(formatCode, wui.KeyControl, wui.KeyF)
	window.SetShortcut(exportPDF, wui.KeyControl, wui.KeyE)
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
			Condition: x.Condition,
			Then:      x.Then,
			TrueText:  x.TrueText,
		}, width, height)

	case parser.IfElse:
		// TODO Paint TrueText and FalseText.
		thenW, thenH := minSize(p, x.Then)
		elseW, elseH := minSize(p, x.Else)
		blockH := max(thenH, elseH)
		// bottom is for the separating line between the condition at the top
		// and the blocks below.
		bottom := height - blockH - 1

		// The width of our available area might be greater than we need them to
		// be for for the two blocks. In that case we want to split the
		// available width in the same ratio as thenW / elseW.
		thenW = int(float64(thenW)/float64(thenW+elseW)*float64(width-1) + 0.5)
		elseW = width - 1 - thenW

		p.Line(0, bottom, width-1, bottom)     // Separate top from bottom.
		p.Line(thenW, bottom, thenW, height-1) // Separate left from right block.
		p.Line(0, 0, thenW, bottom-1)          // Diagonal from top-left.
		p.Line(thenW, bottom-1, width-1, 0)    // Diagonal from top-right.

		_, textH := p.TextSize(x.Condition.Text)
		textH = max(textH, margin)
		// We place the text at the top and offset it from left so that it
		// aligns with the diagonal anchored at the top-left. This diagonal has
		// a slope of textH / y as you can see in the sketch below. We multiply
		// it by thenW to get the right ratio for the text to offset relative to
		// the available width.
		//
		// 	____________________   0
		//       --__  | text | _-
		// 	      --|______|-
		// 	          --__-
		// 	|-------------|-----|  bottom
		//         thenW     elseW
		textX := thenW * textH / bottom
		p.Text(textX+margin/4, 0, x.Condition.Text)

		paintIn(
			offsetPainter{p: p, dy: bottom + 1},
			x.Then,
			thenW, blockH,
		)
		paintIn(
			offsetPainter{p: p, dx: thenW + 1, dy: bottom + 1},
			x.Else,
			elseW, blockH,
		)

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
		textW, textH := p.TextSize(x.Condition.Text)
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
			Condition: x.Condition,
			Then:      x.Then,
			TrueText:  x.TrueText,
		})

	case parser.IfElse:
		// TODO Consider TrueText and FalseText for size.
		thenW, thenH := minSize(p, x.Then)
		elseW, elseH := minSize(p, x.Else)
		textW, textH := p.TextSize(x.Condition.Text)
		textW += margin / 2
		textH = max(textH, margin)
		bottom := max(thenW+1+elseW, textW+textW/2)
		h := int(float64(bottom*textH)/float64(bottom-textW) + 0.5)
		return bottom, h + 1 + max(thenH, elseH)

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
	_, textH := p.TextSize(while.Condition.Text)
	top := margin + textH
	p.Text(margin/2, margin/2, while.Condition.Text)
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
	_, textH := p.TextSize(do.Condition.Text)
	bottom := height - 1 - margin - textH
	p.Line(margin, 0, margin, bottom)
	p.Line(margin, bottom, width-1, bottom)
	p.Text(margin/2, bottom+1+margin/2, do.Condition.Text)
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

type imagePainter struct {
	img  *image.RGBA
	font *gofont.Font
}

var black = color.RGBA{A: 255}

func (p imagePainter) Text(x, y int, s string) {
	p.font.Write(p.img, s, x, y)
}

func (p imagePainter) TextSize(s string) (width, height int) {
	return p.font.Measure(s)
}

func (p imagePainter) Rect(x, y, width, height int) {
	p.Line(x, y, x+width-1, y)
	p.Line(x+width-1, y, x+width-1, y+height-1)
	p.Line(x+width-1, y+height-1, x, y+height-1)
	p.Line(x, y+height-1, x, y)
}

func (p imagePainter) Line(x1, y1, x2, y2 int) {
	// Bresenham's algorithm copied from:
	// http://rosettacode.org/wiki/Bitmap/Bresenham%27s_line_algorithm#Go
	dx := x2 - x1
	if dx < 0 {
		dx = -dx
	}
	dy := y2 - y1
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}
	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		p.img.SetRGBA(x1, y1, black)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func (p imagePainter) LineHeight() int {
	return p.font.HeightInPixels
}
