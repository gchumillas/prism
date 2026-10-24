package ui

import (
	"hash/maphash"
	"image"
	"reflect"
	"slices"
)

func NewHStack(c *Canvas, opts ...StackOption) *Stack {
	return newStack(c, horizontalFootprint, horizontalLayout, opts...)
}

func horizontalFootprint(s *Stack) []uint64 {
	return []uint64{maphash.String(seed, "horizontal"), uint64(s.paddingX), uint64(s.paddingY), uint64(s.gap)}
}

func horizontalLayout(s *Stack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	size = image.Pt(s.paddingX, 0)
	for i, child := range s.children {
		if i > 0 {
			size.X += s.gap
		}
		pos := image.Pt(x+size.X, y+s.paddingY)
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
				prev.regs = MergeRegions(append(prev.regs, childRegs...))
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
		size.X += childSize.X
		size.Y = max(size.Y, childSize.Y)
	}
	for _, child := range s.prevChildStates[len(s.children):] {
		dirtyRegs = append(dirtyRegs, child.regs...)
	}
	clear(s.prevChildStates[len(s.children):])
	s.prevChildStates = s.prevChildStates[:len(s.children)]
	dirtyRegs = MergeRegions(dirtyRegs)
	size = image.Pt(size.X+s.paddingX, size.Y+2*s.paddingY)
	return
}
