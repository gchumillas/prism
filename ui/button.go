package ui

import (
	"image/color"
	"math"
)

type LabelButton struct {
	BasePosition
	scene *Scene
	label *Label
}

func (s *Scene) LabelButton(text string) *LabelButton {
	return &LabelButton{
		scene: s,
		label: s.Label(text),
	}
}

func (lb *LabelButton) Label() *Label {
	return lb.label
}

func (lb *LabelButton) Size() (int, int) {
	w, h := lb.label.Size()
	px, py := lb.padding()
	return w + 2*px, h + 2*py
}

func (lb *LabelButton) Render() {
	w, h := lb.Size()
	px, py := lb.padding()
	lb.label.SetPosition(lb.X+px, lb.Y+py)
	drawRect(lb.scene.canvas, lb.X, lb.Y, lb.X+w, lb.Y+h, color.Black)
	lb.label.Render()
}

func (lb *LabelButton) Contains(x, y int) bool {
	w, h := lb.Size()
	return (lb.X <= x && x <= lb.X+w) && (lb.Y <= y && y <= lb.Y+h)
}

func (lb *LabelButton) padding() (px, py int) {
	fs := lb.label.fontSize
	scale := lb.scene.dpi / 72.0
	var x, y float64
	switch {
	case fs >= 24:
		x, y = 12, 32
	case fs >= 18:
		x, y = 10, 24
	case fs >= 14:
		x, y = 8, 12
	default:
		x, y = 12, 6
	}
	x, y = math.Ceil(scale*x), math.Ceil(scale*y)
	return int(x), int(y)
}
