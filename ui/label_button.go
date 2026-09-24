package ui

var paddingBtn = Point{38, 19}

type LabelButton struct {
	BasePosition
	scene  *Scene
	label  *Label
	hitbox Rectangle
}

func (s *Scene) LabelButton(text string) *LabelButton {
	lb := &LabelButton{
		scene: s,
		label: s.Label(text),
	}
	return lb
}

func (lb *LabelButton) Height() float64 {
	return lb.label.Height() + 2*paddingBtn.Y
}

func (lb *LabelButton) Render() Rectangle {
	width := lb.label.Width() + 2*paddingBtn.X
	lb.label.position = lb.position.Add(paddingBtn)
	lr := lb.label.Render()
	lb.hitbox = Rectangle{
		Min: lr.Min.Sub(paddingBtn),
		Max: Point{max(lr.Max.X+paddingBtn.X, width), lr.Max.Y + paddingBtn.Y},
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
