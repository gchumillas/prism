package ui

import (
	"image"
	"reflect"
	"slices"
)

func NewVStack(opts ...StackOption) *Stack {
	return newStack(verticalLayout, opts...)
}

func verticalLayout(stx *stackContext, st *Stack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	size.Y += st.paddingY
	for i, child := range st.children {
		if i > 0 {
			size.Y += st.gap
		}
		switch child := child.(type) {
		case *Stack:
			childRegs, childSize := child.layout(stx, child, x+st.paddingX, y+size.Y)
			dirtyRegs = append(dirtyRegs, childRegs...)
			size.X, size.Y = max(size.X, childSize.X), size.Y+childSize.Y
		case Widget:
			pos := image.Pt(x+st.paddingX, y+size.Y)
			childRegs, childSize := child.Render(pos.X, pos.Y)
			size.X, size.Y = max(size.X, childSize.X), size.Y+childSize.Y
			footprint := child.Footprint()
			typ := reflect.TypeOf(child)
			if stx.prevChildStatesIdx < len(stx.prevChildStates) {
				prev := &stx.prevChildStates[stx.prevChildStatesIdx]
				hasChanged := prev.pos != pos || prev.typ != typ || !slices.Equal(prev.footprint, footprint)
				if hasChanged {
					prev.pos = pos
					prev.typ = typ
					prev.footprint = footprint
					dirtyRegs = append(dirtyRegs, prev.regs...)
					dirtyRegs = append(dirtyRegs, childRegs...)
					prev.regs = childRegs
				}
			} else {
				stx.prevChildStates = append(stx.prevChildStates, childState{
					pos:       pos,
					typ:       typ,
					footprint: footprint,
					regs:      childRegs,
				})
				dirtyRegs = append(dirtyRegs, childRegs...)
			}
			stx.prevChildStatesIdx++
		default:
			panic("unsupported type")
		}
	}
	size.X += 2 * st.paddingX
	size.Y += st.paddingY
	st.outdated = false
	return
}
