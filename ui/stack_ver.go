package ui

import (
	"hash/maphash"
	"image"
	"reflect"
	"slices"
)

type VStack struct {
	*baseStack
	wrapper *vStackWrapper
}

// TODO: remove param `ctx Context`, since we just need to pass the canvas (affects other places)
func NewVStack(ctx Context, opts ...StackOption) *VStack {
	return &VStack{
		baseStack: newStack(ctx, verticalFootprint, verticalLayout, opts...),
	}
}

type vStackWrapper struct {
	*VStack
}

func (hsw *vStackWrapper) Footprint() []uint64 {
	return hsw.footprint()
}

func AttachVStack(owner Stack, ctx Context, opts ...StackOption) *VStack {
	vs := NewVStack(ctx, opts...)
	vs.wrapper = &vStackWrapper{VStack: vs}
	owner.Add(vs.wrapper)
	return vs
}

func (s *VStack) DetachFrom(owner Stack) bool {
	return owner.Remove(s.wrapper) > 0
}

func verticalFootprint(s *baseStack) []uint64 {
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
	maphash.WriteComparable(&h, s.padding.X)
	maphash.WriteComparable(&h, s.padding.Y)
	maphash.WriteComparable(&h, s.gap)
	return append([]uint64{h.Sum64()}, result...)
}

func verticalLayout(s *baseStack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
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
