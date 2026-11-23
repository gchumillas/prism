package ui

import (
	"image"
	"reflect"
	"slices"
)

func NewHStack(opts ...StackOption) *Stack {
	return newStack(horizontalLayout, horizontalLayoutSize, opts...)
}

func horizontalLayoutSize(st *Stack) (size image.Point) {
	size.X = st.paddingX
	for i, child := range st.children {
		if i > 0 {
			size.X += st.gap
		}
		switch child := child.(type) {
		case *Stack:
			childSize := child.layoutSize(child)
			size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
		case Widget:
			childSize := child.Size()
			size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
		default:
			panic("unsupported type")
		}
	}
	size.X += st.paddingX
	size.Y += 2 * st.paddingY
	return
}

func horizontalLayout(stx *stackContext, st *Stack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	childStatesStartIdx := stx.prevChildStatesIdx
	size.X = st.paddingX
	for i, child := range st.children {
		if i > 0 {
			size.X += st.gap
		}
		pos := image.Pt(x+size.X, y+st.paddingY)
		switch child := child.(type) {
		case *Stack:
			childRegs, childSize := child.layout(stx, child, pos.X, pos.Y)
			dirtyRegs = append(dirtyRegs, childRegs...)
			size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
		case Widget:
			childRegs := child.Render(pos.X, pos.Y)
			childSize := child.Size()
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
	size.X += st.paddingX
	size.Y += 2 * st.paddingY
	st.prevChildStates = slices.Clone(stx.prevChildStates[childStatesStartIdx:stx.prevChildStatesIdx])
	st.outdated = false
	return
}
