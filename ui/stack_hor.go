package ui

import (
	"image"
	"reflect"
	"slices"
)

func NewHStack() *Stack {
	return newStack(horizontalLayout, horizontalLayoutSize)
}

func horizontalLayoutSize(st *Stack) (size image.Point) {
	size.X = st.PaddingX
	for i, child := range st.children {
		if i > 0 {
			size.X += st.Gap
		}
		childSize := child.Size()
		size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
	}
	size.X += st.PaddingX
	size.Y += 2 * st.PaddingY
	return
}

func horizontalLayout(stx *stackContext, st *Stack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	childStatesStartIdx := stx.prevChildStatesIdx
	size.X = st.PaddingX
	for i, child := range st.children {
		if i > 0 {
			size.X += st.Gap
		}
		pos := image.Pt(x+size.X, y+st.PaddingY)
		switch v := child.(type) {
		case *Stack:
			childRegs, childSize := v.layout(stx, v, pos.X, pos.Y)
			dirtyRegs = append(dirtyRegs, childRegs...)
			size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
		case Widget:
			childRegs := child.Render(pos.X, pos.Y)
			childSize := child.Size()
			size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
			typ := reflect.TypeOf(child)
			if stx.prevChildStatesIdx < len(stx.prevChildStates) {
				footprint := v.Footprint()
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
					footprint: v.Footprint(),
					regs:      childRegs,
				})
				dirtyRegs = append(dirtyRegs, childRegs...)
			}
			stx.prevChildStatesIdx++
		default:
			panic("unsupported type")
		}
	}
	size.X += st.PaddingX
	size.Y += 2 * st.PaddingY
	st.prevChildStates = slices.Clone(stx.prevChildStates[childStatesStartIdx:stx.prevChildStatesIdx])
	st.outdated = false
	return
}
