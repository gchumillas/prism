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
		childSize := child.Size()
		size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
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
		if v, ok := child.(*Stack); ok {
			childRegs, childSize := v.layout(stx, v, pos.X, pos.Y)
			dirtyRegs = append(dirtyRegs, childRegs...)
			size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
		} else {
			childRegs := child.Render(pos.X, pos.Y)
			childSize := child.Size()
			size.X, size.Y = size.X+childSize.X, max(size.Y, childSize.Y)
			var footprint []uint64
			switch v := child.(type) {
			case *Button:
				footprint = v.widget.Footprint()
			case Widget:
				footprint = v.Footprint()
			default:
				panic("unsupported type")
			}
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
		}
	}
	size.X += st.paddingX
	size.Y += 2 * st.paddingY
	st.prevChildStates = slices.Clone(stx.prevChildStates[childStatesStartIdx:stx.prevChildStatesIdx])
	st.outdated = false
	return
}
