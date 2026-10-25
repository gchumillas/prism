package ui

import (
	"image"
	"reflect"
	"slices"
)

func NewVStack(c *Canvas, opts ...StackOption) *Stack {
	return newStack(c, verticalLayout, opts...)
}

func verticalLayout(stx *stackContext, st *Stack, x, y int) {
	stx.size.Y += st.paddingY
	for i, child := range st.children {
		if i > 0 {
			stx.size.Y += st.gap
		}
		switch child := child.(type) {
		case *Stack:
			prevSize := stx.size
			stx.size = image.Point{}
			child.layout(stx, child, x+st.paddingX, y+prevSize.Y)
			stx.size.X = max(prevSize.X, stx.size.X)
			stx.size.Y += prevSize.Y
		case Widget:
			pos := image.Pt(x+st.paddingX, y+stx.size.Y)
			childRegs, childSize := child.Render(pos.X, pos.Y)
			stx.size.X = max(stx.size.X, childSize.X)
			stx.size.Y += childSize.Y
			footprint := child.Footprint()
			if stx.prevChildStatesIdx < len(stx.prevChildStates) {
				prev := &stx.prevChildStates[stx.prevChildStatesIdx]
				hasChanged := prev.pos != pos || prev.typ != reflect.TypeOf(child) || !slices.Equal(prev.footprint, footprint)
				if hasChanged {
					prev.pos = pos
					prev.typ = reflect.TypeOf(child)
					prev.footprint = footprint
					stx.dirtyRegs = append(stx.dirtyRegs, prev.regs...)
					stx.dirtyRegs = append(stx.dirtyRegs, childRegs...)
					prev.regs = childRegs
				}
			} else {
				stx.prevChildStates = append(stx.prevChildStates, childState{
					pos:       pos,
					typ:       reflect.TypeOf(child),
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
	stx.size.X += 2 * st.paddingX
	stx.size.Y += st.paddingY
	st.outdated = false
}
