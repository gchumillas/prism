package ui

import (
	"image"
	"reflect"
	"slices"
)

func NewHStack(opts ...StackOption) *Stack {
	return newStack(horizontalLayout, opts...)
}

func horizontalLayout(stx *stackContext, st *Stack, x, y int) {
	stx.size.X += st.paddingX
	for i, child := range st.children {
		if i > 0 {
			stx.size.X += st.gap
		}
		switch child := child.(type) {
		case *Stack:
			prevSize := stx.size
			stx.size = image.Point{}
			child.layout(stx, child, x+prevSize.X, y+st.paddingY)
			stx.size.X += prevSize.X
			stx.size.Y = max(prevSize.Y, stx.size.Y)
		case Widget:
			pos := image.Pt(x+stx.size.X, y+st.paddingY)
			childRegs, childSize := child.Render(pos.X, pos.Y)
			stx.size.X += childSize.X
			stx.size.Y = max(stx.size.Y, childSize.Y)
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
	stx.size.X += st.paddingX
	stx.size.Y += 2 * st.paddingY
	st.outdated = false
}
