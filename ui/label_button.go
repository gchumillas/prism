package ui

var paddingBtn = Point{38, 19}

type LabelButton struct {
	scene  *Scene
	hitbox Rectangle
	vstack *VStack
	label  *Label
}

func (s *Scene) LabelButton(text string) *LabelButton {
	lb := &LabelButton{scene: s}
	lb.vstack = s.VStack(Padding(paddingBtn.X, paddingBtn.Y))
	lb.label = s.Label(text)
	lb.vstack.Add(lb.label)
	return lb
}

func (lb *LabelButton) X() float64       { return lb.vstack.X() }
func (lb *LabelButton) SetX(val float64) { lb.vstack.SetX(val) }

func (lb *LabelButton) Y() float64       { return lb.vstack.Y() }
func (lb *LabelButton) SetY(val float64) { lb.vstack.SetY(val) }

func (lb *LabelButton) Height() float64 {
	return lb.label.Height() + 2*paddingBtn.Y
}

func (lb *LabelButton) Render() Rectangle {
	size := lb.label.Size().Add(paddingBtn.Mul(2))
	lb.hitbox = Rectangle{
		Min: lb.vstack.position,
		Max: lb.vstack.position.Add(size),
	}

	lb.scene.canvas.SetLineWidth(4)
	lb.scene.canvas.DrawRectangle(lb.hitbox.Min.X, lb.hitbox.Min.Y, lb.hitbox.Max.X-lb.hitbox.Min.X, lb.hitbox.Max.Y-lb.hitbox.Min.Y)
	lb.scene.canvas.Stroke()

	return lb.hitbox
}

func (lb *LabelButton) Contains(x, y float64) bool {
	p := Point{x, y}
	return p.In(lb.hitbox)
}
