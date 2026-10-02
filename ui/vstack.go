package ui

import (
	"image"
	"slices"
)

type VStackChild interface {
	Renderer
	Height() int
}

type VStack struct {
	Padding image.Point
	Gap     int

	children []VStackChild
}

func NewVStack(c Context) *VStack {
	return &VStack{}
}

func (vs *VStack) Height() int {
	height := vs.Padding.Y
	for i, child := range vs.children {
		if i > 0 {
			height += vs.Gap
		}
		height += child.Height()
	}
	return height
}

func (vs *VStack) Children() []VStackChild {
	return vs.children
}

func (vs *VStack) Add(child VStackChild) {
	vs.children = append(vs.children, child)
}

func (vs *VStack) Remove(child VStackChild) {
	vs.children = slices.DeleteFunc(vs.children, func(c VStackChild) bool {
		return c == child
	})
}

func (vs *VStack) Render(pos image.Point) (dirtyRegs []image.Rectangle) {
	y := pos.Y + vs.Padding.Y
	for i, child := range vs.children {
		if i > 0 {
			y += vs.Gap + vs.children[i-1].Height()
		}
		dirtyRegs = append(dirtyRegs, child.Render(image.Pt(pos.X+vs.Padding.X, y))...)
	}
	return
}
