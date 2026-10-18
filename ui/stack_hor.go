package ui

import (
	"hash/maphash"
	"image"
	"reflect"
	"slices"
)

type HStack struct {
	*baseStack
	wrapper *hStackWrapper
}

func NewHStack(ctx Context, opts ...StackOption) *HStack {
	return &HStack{
		baseStack: newStack(ctx, horizontalFootprint, horizontalLayout, opts...),
	}
}

type hStackWrapper struct {
	*HStack
}

func (hsw *hStackWrapper) Footprint() []uint64 {
	return hsw.footprint()
}

func AttachHStack(owner Stack, ctx Context, opts ...StackOption) *HStack {
	hs := NewHStack(ctx, opts...)
	hs.wrapper = &hStackWrapper{HStack: hs}
	owner.Add(hs.wrapper)
	return hs
}

func (s *HStack) DetachFrom(owner Stack) bool {
	return owner.Remove(s.wrapper) > 0
}

func horizontalFootprint(s *baseStack) []uint64 {
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
	maphash.WriteComparable(&h, s.padding.X)
	maphash.WriteComparable(&h, s.padding.Y)
	maphash.WriteComparable(&h, s.gap)
	return append([]uint64{h.Sum64()}, result...)
}

func horizontalLayout(s *baseStack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	size = image.Pt(s.padding.X, 0)
	for i, child := range s.children {
		if i > 0 {
			size.X += s.gap
		}
		pos := image.Pt(x+size.X, y+s.padding.Y)
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
	size = image.Pt(size.X+s.padding.X, size.Y+2*s.padding.Y)
	return
}
