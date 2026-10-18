package ui

import (
	"hash/maphash"
	"image"
	"reflect"
	"slices"
)

func NewVStack(c *Canvas, opts ...stackOption) *Stack {
	return newStack(c, verticalFootprint, verticalLayout, opts...)
}

func AttachVStackTo(parent *Stack, c *Canvas, opts ...stackOption) *Stack {
	vs := NewVStack(c, opts...)
	vs.wrapper = &stackWrapper{Stack: vs}
	parent.Add(vs.wrapper)
	return vs
}

func verticalFootprint(s *Stack) []uint64 {
	var result []uint64
	for _, child := range s.children {
		typ := typeID(reflect.TypeOf(child))
		footprint := child.Footprint()
		result = append(result, uint64(typ), uint64(len(footprint)))
		result = append(result, footprint...)
	}

	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, "vertical")
	maphash.WriteComparable(&h, s.paddingX)
	maphash.WriteComparable(&h, s.paddingY)
	maphash.WriteComparable(&h, s.gap)
	return append([]uint64{h.Sum64()}, result...)
}

func verticalLayout(s *Stack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	childPos := image.Pt(pos.X+s.paddingX, pos.Y+s.paddingY)
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
