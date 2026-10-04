package ui

import (
	"bytes"
	"image"
	"slices"
)

type VStack struct {
	Padding  image.Point
	Gap      int
	children []Widget

	prevPos         image.Point
	prevChildStates map[int][]byte
}

func NewVStack(c Context) *VStack {
	return &VStack{
		prevChildStates: map[int][]byte{},
	}
}

func (vs *VStack) Size() image.Point {
	width, height := 0, vs.Padding.Y
	for i, child := range vs.children {
		if i > 0 {
			height += vs.Gap
		}
		size := child.Size()
		width = max(width, size.X)
		height += size.Y
	}
	return image.Pt(width+vs.Padding.X, height)
}

func (vs *VStack) Children() []Widget {
	return vs.children
}

func (vs *VStack) Add(child Widget) {
	vs.children = append(vs.children, child)
}

func (vs *VStack) Remove(child Widget) {
	vs.children = slices.DeleteFunc(vs.children, func(c Widget) bool {
		return c == child
	})
}

func (vs *VStack) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	y0 := pos.Y + vs.Padding.Y
	for i, child := range vs.children {
		if i > 0 {
			y0 += vs.Gap + vs.children[i-1].Size().Y
		}
		state := serializeNode(x, y0, child)
		prevState, ok := vs.prevChildStates[i]
		if vs.prevPos != pos || !ok || !bytes.Equal(prevState, state) {
			dirtyRegs = append(dirtyRegs, child.Render(pos.X+vs.Padding.X, y0)...)
		}
		vs.prevChildStates[i] = state
	}
	vs.prevPos = pos
	return
}
