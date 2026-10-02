package ui

import (
	"image"
	"slices"
)

type VStackChild interface {
	Renderer
	SetPosition(image.Point)
	Height() int
}

type VStack struct {
	Position image.Point
	Padding  image.Point
	Gap      int

	children []VStackChild
}

func NewVStack(c Context) *VStack {
	return &VStack{}
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

func (vs *VStack) Render() (dirtyRegs []image.Rectangle) {
	height := vs.Position.Y + vs.Padding.Y
	for i, child := range vs.children {
		if i > 0 {
			height += vs.Gap + vs.children[i-1].Height()
		}
		child.SetPosition(image.Pt(vs.Position.X+vs.Padding.X, height))
		dirtyRegs = append(dirtyRegs, child.Render()...)
	}
	return
}
