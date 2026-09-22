package ui

import (
	"fmt"
	"image"
)

type LayoutChild interface {
	Renderer
	SetPosition(x, y int)
	SetWidth(int)
}

type VStack struct {
	x, y     int
	width    int
	gap      int
	children []LayoutChild
}

func (s *Scene) VStack(width int, opts ...VStackOpt) *VStack {
	vs := &VStack{width: width}
	for _, opt := range opts {
		opt.applyVStack(vs)
	}
	return vs
}

func (vs *VStack) Add(child LayoutChild) {
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
		child.SetPosition(vs.x, y)
		child.SetWidth(vs.width)
		r := child.Render()
		fmt.Println(r)
		y += vs.gap + (r.Max.Y - r.Min.Y)
	}
	return
}
