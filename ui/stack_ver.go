package ui

import (
	"image"
	"reflect"
	"slices"
)

func NewVStack(opts ...StackOption) *Stack {
	return newStack(verticalLayout, verticalLayoutSize, opts...)
}

func verticalLayoutSize(st *Stack) (size image.Point) {
	size.Y += st.paddingY
	for i, child := range st.children {
		if i > 0 {
			size.Y += st.gap
		}
		childSize := child.Size()
		size.X, size.Y = max(size.X, childSize.X), size.Y+childSize.Y
	}
	size.X += 2 * st.paddingX
	size.Y += st.paddingY
	return
}

func verticalLayout(stx *stackContext, st *Stack, x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	childStatesStartIdx := stx.prevChildStatesIdx
	size.Y += st.paddingY
	for i, child := range st.children {
		if i > 0 {
			size.Y += st.gap
		}
		pos := image.Pt(x+st.paddingX, y+size.Y)
		switch v := child.(type) {
		case *Stack:
			childRegs, childSize := v.layout(stx, v, pos.X, pos.Y)
			dirtyRegs = append(dirtyRegs, childRegs...)
			size.X, size.Y = max(size.X, childSize.X), size.Y+childSize.Y
		case Widget:
			childRegs := v.Render(pos.X, pos.Y)
			childSize := v.Size()
			size.X, size.Y = max(size.X, childSize.X), size.Y+childSize.Y
			typ := reflect.TypeOf(v)
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
	size.X += 2 * st.paddingX
	size.Y += st.paddingY
	st.prevChildStates = slices.Clone(stx.prevChildStates[childStatesStartIdx:stx.prevChildStatesIdx])
	st.outdated = false
	return
}
