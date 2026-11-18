package ui

import (
	"image"
	"reflect"
	"slices"
)

func NewHStack(opts ...StackOption) *Stack {
	return newStack(horizontalLayout, opts...)
}

func horizontalLayout(stx *stackContext, st *Stack, x, y int) (size image.Point) {
	size.X = st.paddingX
	for i, child := range st.children {
		if i > 0 {
			size.X += st.gap
		}
		switch child := child.(type) {
		case *Stack:
			childSize := child.layout(stx, child, x+size.X, y+st.paddingY)
			size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
		case Widget:
			pos := image.Pt(x+size.X, y+st.paddingY)
			childRegs, childSize := child.Render(pos.X, pos.Y)
			size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
			footprint := child.Footprint()
			typ := reflect.TypeOf(child)
			if stx.prevChildStatesIdx < len(stx.prevChildStates) {
				prev := &stx.prevChildStates[stx.prevChildStatesIdx]
				hasChanged := prev.pos != pos || prev.typ != typ || !slices.Equal(prev.footprint, footprint)
				if hasChanged {
					prev.pos = pos
					prev.typ = typ
					prev.footprint = footprint
					stx.dirtyRegs = append(stx.dirtyRegs, prev.regs...)
					stx.dirtyRegs = append(stx.dirtyRegs, childRegs...)
					prev.regs = childRegs
				}
			} else {
				stx.prevChildStates = append(stx.prevChildStates, childState{
					pos:       pos,
					typ:       typ,
					footprint: footprint,
					regs:      childRegs,
				})
				stx.dirtyRegs = append(stx.dirtyRegs, childRegs...)
			}
			stx.prevChildStatesIdx++
		default:
			panic("unsupported type")
		}
	}
	size.X += st.paddingX
	size.Y += 2 * st.paddingY
	st.outdated = false
	return
}
