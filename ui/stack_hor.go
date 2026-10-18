package ui

import (
	"hash/maphash"
	"image"
	"reflect"
	"slices"
)

func NewHStack(c *Canvas, opts ...stackOption) *Stack {
	return newStack(c, horizontalFootprint, horizontalLayout, opts...)
}

func AttachHStackTo(parent *Stack, c *Canvas, opts ...stackOption) *Stack {
	hs := NewHStack(c, opts...)
	hs.wrapper = &stackWrapper{Stack: hs}
	parent.Add(hs.wrapper)
	return hs
}

func horizontalFootprint(s *Stack) []uint64 {
	var result []uint64
	for _, child := range s.children {
		typ := typeID(reflect.TypeOf(child))
		footprint := child.Footprint()
		result = append(result, uint64(typ), uint64(len(footprint)))
		result = append(result, footprint...)
	}

	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, "horizontal")
	maphash.WriteComparable(&h, s.paddingX)
	maphash.WriteComparable(&h, s.paddingY)
	maphash.WriteComparable(&h, s.gap)
	return append([]uint64{h.Sum64()}, result...)
}

func horizontalLayout(s *Stack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	size = image.Pt(s.paddingX, 0)
	for i, child := range s.children {
		if i > 0 {
			size.X += s.gap
		}
		pos := image.Pt(x+size.X, y+s.paddingY)
		childRegs, childSize := child.Render(pos.X, pos.Y)
		cs := childState{
			pos:       pos,
			typ:       reflect.TypeOf(child),
			footprint: child.Footprint(),
			// TODO: this is wrong, since "size" and "dirty regions" are different concepts
			regs: []image.Rectangle{{Min: pos, Max: pos.Add(childSize)}},
		}
		if i < len(s.prevChildStates) {
			prev := &s.prevChildStates[i]
			if prev.pos != cs.pos || prev.typ != cs.typ || !slices.Equal(prev.footprint, cs.footprint) {
				dirtyRegs = append(dirtyRegs, prev.regs...)
				dirtyRegs = append(dirtyRegs, childRegs...)
				*prev = cs
			}
		} else {
			s.prevChildStates = append(s.prevChildStates, cs)
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
	size = image.Pt(size.X+s.paddingX, size.Y+2*s.paddingY)
	return
}
