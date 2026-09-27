package ui

import "image"

type VStackChild interface {
	Renderer
	SetX(int)
	SetY(int)
	Height() int
}

type VStack struct {
	BasePosition
	padding  image.Point
	gap      int
	children []VStackChild
}

type VStackOpt interface {
	applyVStack(*VStack)
}

// avail options
func (val paddingXOpt) applyVStack(x *VStack) { x.padding.X = int(val) }
func (val paddingYOpt) applyVStack(x *VStack) { x.padding.Y = int(val) }
func (val gapOpt) applyVStack(x *VStack)      { x.gap = int(val) }

func (c Context) VStack(opts ...VStackOpt) *VStack {
	vs := &VStack{}
	for _, opt := range opts {
		opt.applyVStack(vs)
	}
	return vs
}

func (vs *VStack) Add(child VStackChild) {
	vs.children = append(vs.children, child)
}

func (vs *VStack) Render() (r image.Rectangle) {
	height := vs.position.Y + vs.padding.Y
	for i, child := range vs.children {
		if i > 0 {
			height += vs.gap + vs.children[i-1].Height()
		}
		child.SetX(vs.position.X + vs.padding.X)
		child.SetY(height)
		r = r.Union(child.Render())
	}
	return
}
