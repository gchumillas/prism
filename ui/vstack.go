package ui

import (
	"image"
)

type VStackChild interface {
	Renderer
	SetX(int)
	SetY(int)
	SetWidth(int)
}

type VStack struct {
	BasePosition
	width    int
	gap      int
	children []VStackChild
}

func (s *Scene) VStack(width int, opts ...VStackOpt) *VStack {
	vs := &VStack{width: width}
	for _, opt := range opts {
		opt.applyVStack(vs)
	}
	return vs
}

func (vs *VStack) Add(child VStackChild) {
	vs.children = append(vs.children, child)
}

func (vs *VStack) Render() (r image.Rectangle) {
	y := vs.position.Y
	for _, child := range vs.children {
		child.SetX(vs.position.X)
		child.SetY(y)
		child.SetWidth(vs.width)
		r := child.Render()
		y += vs.gap + (r.Max.Y - r.Min.Y)
	}
	return
}
