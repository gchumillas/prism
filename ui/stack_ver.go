package ui

import (
	"hash/maphash"
	"image"
	"reflect"
	"slices"
)

func NewVStack(c *Canvas, opts ...StackOption) *Stack {
	return newStack(c, verticalFootprint, verticalLayout, opts...)
}

func verticalFootprint(s *Stack) []uint64 {
	return []uint64{maphash.String(seed, "vertical"), uint64(s.paddingX), uint64(s.paddingY), uint64(s.gap)}
}

func verticalLayout(s *Stack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	size = image.Pt(0, s.paddingY)
	for i, child := range s.children {
		if i > 0 {
			size.Y += s.gap
		}
		pos := image.Pt(x+s.paddingX, y+size.Y)
		childRegs, childSize := child.Render(pos.X, pos.Y)
		footprint := child.Footprint()
		if i < len(s.prevChildStates) {
			prev := &s.prevChildStates[i]
			hasChanged := prev.pos != pos || prev.typ != reflect.TypeOf(child) || !slices.Equal(prev.footprint, footprint)
			if hasChanged {
				prev.pos = pos
				prev.typ = reflect.TypeOf(child)
				prev.footprint = footprint
			}
			if _, ok := child.(*stackWrapper); ok {
				dirtyRegs = append(dirtyRegs, childRegs...)
				if hasChanged {
					dirtyRegs = append(dirtyRegs, prev.regs...)
				}
				if len(childRegs) > 0 {
					prev.regs = MergeRegions(append(prev.regs, childRegs...))
				}
			} else if hasChanged {
				dirtyRegs = append(dirtyRegs, prev.regs...)
				dirtyRegs = append(dirtyRegs, childRegs...)
				prev.regs = childRegs
			}
		} else {
			s.prevChildStates = append(s.prevChildStates, childState{
				pos:       pos,
				typ:       reflect.TypeOf(child),
				footprint: footprint,
				regs:      childRegs,
			})
			dirtyRegs = append(dirtyRegs, childRegs...)
		}
		size.X = max(size.X, childSize.X)
		size.Y += childSize.Y
	}
	for _, child := range s.prevChildStates[len(s.children):] {
		dirtyRegs = append(dirtyRegs, child.regs...)
	}
	clear(s.prevChildStates[len(s.children):])
	s.prevChildStates = s.prevChildStates[:len(s.children)]
	dirtyRegs = MergeRegions(dirtyRegs)
	size = image.Pt(size.X+2*s.paddingX, size.Y+s.paddingY)
	return
}
