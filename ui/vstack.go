package ui

import "image"

type VStackChild interface {
	Renderer
	// TODO: replace with SetPosition
	SetX(int)
	SetY(int)
	// TODO: replace with Size()?
	Height() int
}

type VStack struct {
	Position image.Point
	Padding  image.Point
	Gap      int
	children []VStackChild
}

func (c Context) VStack() *VStack {
	return &VStack{}
}

func (vs *VStack) Add(child VStackChild) {
	vs.children = append(vs.children, child)
}

func (vs *VStack) Render() (dirtyRegs []image.Rectangle) {
	height := vs.Position.Y + vs.Padding.Y
	for i, child := range vs.children {
		if i > 0 {
			height += vs.Gap + vs.children[i-1].Height()
		}
		child.SetX(vs.Position.X + vs.Padding.X)
		child.SetY(height)
		dirtyRegs = append(dirtyRegs, child.Render()...)
	}
	return
}
