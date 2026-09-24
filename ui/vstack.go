package ui

type VStackChild interface {
	Renderer
	SetX(float64)
	SetY(float64)
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
	y := vs.position.Y
	for _, child := range vs.children {
		child.SetX(vs.position.X)
		child.SetY(y)
		r := child.Render()
		y += vs.gap + (r.Max.Y - r.Min.Y)
	}
	return
}
