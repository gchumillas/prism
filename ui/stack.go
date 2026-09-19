package ui

type VStack struct {
	elems []Element
	x, y  int
	gap   int
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

func (l *VStack) Add(elem Element) {
	l.elems = append(l.elems, elem)
}

func (l *VStack) Size() (int, int) {
	w0, h0 := 0, 0
	for i, e := range l.elems {
		if i > 0 {
			h0 += l.gap
		}
		w, h := e.Size()
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

func (l *VStack) Draw() {
	y := l.y
	for _, e := range l.elems {
		_, h := e.Size()
		e.SetPosition(l.x, y)
		e.Draw()
		y += l.gap + h
	}
}
