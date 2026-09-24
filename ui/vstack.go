package ui

type VStackChild interface {
	Renderer
	SetX(float64)
	SetY(float64)
	Height() float64
}

type VStack struct {
	BasePosition
	gap      float64
	children []VStackChild
}

func (s *Scene) VStack(opts ...VStackOpt) *VStack {
	vs := &VStack{}
	for _, opt := range opts {
		opt.applyVStack(vs)
	}
	return vs
}

func (vs *VStack) Add(child VStackChild) {
	vs.children = append(vs.children, child)
}

func (vs *VStack) Render() (r Rectangle) {
	height := vs.position.Y
	for i, child := range vs.children {
		if i > 0 {
			height += vs.gap + vs.children[i-1].Height()
		}
		child.SetX(vs.position.X)
		child.SetY(height)
		r = r.Union(child.Render())
	}
	return
}
