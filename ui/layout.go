package ui

type LayoutChild interface {
	Renderer
	SetPosition(x, y int)
	Size() (width, height int)
}

type VStack struct {
	x, y     int
	gap      int
	children []LayoutChild
}

func (s *Scene) VStack() *VStack {
	return &VStack{}
}

func (l *VStack) Gap() int {
	return l.gap
}

func (l *VStack) SetGap(val int) {
	l.gap = val
}

func (vs *VStack) Add(child LayoutChild) {
	vs.children = append(vs.children, child)
}

func (vs *VStack) Size() (int, int) {
	w0, h0 := 0, 0
	for i, child := range vs.children {
		if i > 0 {
			h0 += vs.gap
		}
		w, h := child.Size()
		if w > w0 {
			w0 = w
		}
		h0 += h
	}
	return w0, h0
}

func (vs *VStack) Position() (int, int) {
	return vs.x, vs.y
}

func (vs *VStack) SetPosition(x, y int) {
	vs.x, vs.y = x, y
}

func (vs *VStack) Render() {
	y := vs.y
	for _, child := range vs.children {
		_, h := child.Size()
		child.SetPosition(vs.x, y)
		child.Render()
		y += vs.gap + h
	}
}
