package ui

import (
	"image"
	"reflect"
	"slices"
)

type VStack struct {
	stack
}

func NewVStack(ctx Context, opts ...stackOption) *VStack {
	s := &VStack{}
	s.ctx = ctx
	for _, opt := range opts {
		opt.applyStack(&s.stack)
	}
	return s
}

func (s *VStack) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	childPos := image.Pt(pos.X+s.padding.X, pos.Y+s.padding.Y)
	prevHeight := 0
	for i, child := range s.children {
		if i > 0 {
			childPos.Y += s.gap + prevHeight
			size.Y += s.gap
		}
		regs, childSize := child.Render(childPos.X, childPos.Y)
		prevHeight = childSize.Y
		size.X = max(size.X, childSize.X)
		size.Y += childSize.Y
		if i < len(s.prevChildStates) {
			prev := s.prevChildStates[i]
			if prev.pos == childPos && prev.typ == reflect.TypeOf(child) && slices.Equal(prev.footprint, child.Footprint()) {
				continue
			}
		}

		dirtyRegs = append(dirtyRegs, regs...)
		cs := childState{pos: childPos, regs: regs, typ: reflect.TypeOf(child), footprint: child.Footprint()}
		if i < len(s.prevChildStates) {
			dirtyRegs = append(dirtyRegs, s.prevChildStates[i].regs...)
			s.prevChildStates[i] = cs
		} else {
			s.prevChildStates = append(s.prevChildStates, cs)
		}
	}
	for _, cs := range s.prevChildStates[len(s.children):] {
		dirtyRegs = append(dirtyRegs, cs.regs...)
	}
	s.prevChildStates = s.prevChildStates[:len(s.children)]
	dirtyRegs = MergeRegions(dirtyRegs)
	return
}
