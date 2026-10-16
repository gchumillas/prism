package ui

import (
	"image"
	"reflect"
	"slices"
)

type HStack struct {
	stack
}

func NewHStack(ctx Context, opts ...stackOption) *HStack {
	s := &HStack{}
	s.ctx = ctx
	for _, opt := range opts {
		opt.applyStack(&s.stack)
	}
	return s
}

func (s *HStack) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	size = image.Pt(s.padding.X, 0)
	for i, child := range s.children {
		if i > 0 {
			size.X += s.gap
		}
		pos := image.Pt(x+size.X, y+s.padding.Y)
		childRegs, childSize := child.Render(pos.X, pos.Y)
		if i < len(s.prevChildStates) {
			prev := &s.prevChildStates[i]
			if prev.pos != pos || prev.typ != reflect.TypeOf(child) || !slices.Equal(prev.footprint, child.Footprint()) {
				prev.pos = pos
				prev.typ = reflect.TypeOf(child)
				prev.footprint = child.Footprint()
				// MergeRegions() helps reduce the number of rectangles
				prev.regs = MergeRegions(append(prev.regs, childRegs...))
				dirtyRegs = append(dirtyRegs, prev.regs...)
			}
		} else {
			s.prevChildStates = append(s.prevChildStates, childState{
				pos:       pos,
				typ:       reflect.TypeOf(child),
				footprint: child.Footprint(),
				regs:      childRegs,
			})
			dirtyRegs = append(dirtyRegs, childRegs...)
		}
		size.X += childSize.X
		size.Y = max(size.Y, childSize.Y)
	}
	for _, child := range s.prevChildStates[len(s.children):] {
		dirtyRegs = append(dirtyRegs, child.regs...)
	}
	// TODO: button leaks?
	s.prevChildStates = s.prevChildStates[:len(s.children)]
	dirtyRegs = MergeRegions(dirtyRegs)
	size = image.Pt(size.X+s.padding.X, size.Y+2*s.padding.Y)
	return
}
