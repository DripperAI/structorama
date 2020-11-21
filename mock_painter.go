package main

import (
	"fmt"
	"sort"
	"testing"
)

type mockPainter struct {
	lineHeight   int
	textW, textH int
	textSizes    map[string][2]int
	ops          []string
}

func (p *mockPainter) LineHeight() int {
	return p.lineHeight
}

func (p *mockPainter) TextSize(s string) (width, height int) {
	if p.textSizes != nil {
		size := p.textSizes[s]
		return size[0], size[1]
	}
	return p.textW, p.textH
}

func (p *mockPainter) Line(x1, y1, x2, y2 int) {
	p.ops = append(p.ops, fmt.Sprintf("Line(%d, %d, %d, %d)", x1, y1, x2, y2))
}

func (p *mockPainter) Rect(x, y, width, height int) {
	p.ops = append(p.ops, fmt.Sprintf("Rect(%d, %d, %d, %d)", x, y, width, height))
}

func (p *mockPainter) Text(x, y int, s string) {
	p.ops = append(p.ops, fmt.Sprintf("Text(%d, %d, %q)", x, y, s))
}

func (p *mockPainter) checkPainting(t *testing.T, wantOps ...string) {
	t.Helper()

	sort.Strings(wantOps)
	sort.Strings(p.ops)

	for _, s := range wantOps {
		if !contains(p.ops, s) {
			t.Error("Missing:", s)
		}
	}

	for _, s := range p.ops {
		if !contains(wantOps, s) {
			t.Error("  Extra:", s)
		}
	}
}

func contains(list []string, s string) bool {
	for _, s2 := range list {
		if s == s2 {
			return true
		}
	}
	return false
}
