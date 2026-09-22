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
	x, y     int
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

func (vs *VStack) Position() (int, int) {
	return vs.x, vs.y
}

func (vs *VStack) SetPosition(x, y int) {
	vs.x, vs.y = x, y
}

func (vs *VStack) Render() (r image.Rectangle) {
	y := vs.y
	for _, child := range vs.children {
		child.SetX(vs.x)
		child.SetY(y)
		child.SetWidth(vs.width)
		r := child.Render()
		y += vs.gap + (r.Max.Y - r.Min.Y)
	}
	return
}
