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

func (l *VStack) Add(child LayoutChild) {
	l.children = append(l.children, child)
}

func (l *VStack) Size() (int, int) {
	w0, h0 := 0, 0
	for i, child := range l.children {
		if i > 0 {
			h0 += l.gap
		}
		w, h := child.Size()
		if w > w0 {
			w0 = w
		}
		h0 += l.gap + h
	}
	return w0, h0
}

func (l *VStack) Position() (int, int) {
	return l.x, l.y
}

func (l *VStack) SetPosition(x, y int) {
	l.x, l.y = x, y
}

func (l *VStack) Render() {
	y := l.y
	for _, child := range l.children {
		_, h := child.Size()
		child.SetPosition(l.x, y)
		child.Render()
		y += l.gap + h
	}
}
