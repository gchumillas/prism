package ui

import (
	"image"
	"slices"
)

type VStack struct {
	Padding image.Point
	Gap     int

	children []Widget
}

func NewVStack(c Context) *VStack {
	return &VStack{}
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
	h := pos.Y + vs.Padding.Y
	for i, child := range vs.children {
		if i > 0 {
			h += vs.Gap + vs.children[i-1].Size().Y
		}
		dirtyRegs = append(dirtyRegs, child.Render(pos.X+vs.Padding.X, h)...)
	}
	return
}
