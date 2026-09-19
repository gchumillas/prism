package ui

import (
	"image/color"
	"math"
)

type LabelButton struct {
	scene *Scene
	label *Label
	x, y  int
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

func (lb *LabelButton) Position() (int, int) {
	return lb.x, lb.y
}

func (lb *LabelButton) SetPosition(x, y int) {
	lb.x, lb.y = x, y
}

func (lb *LabelButton) Draw() {
	w, h := lb.Size()
	px, py := lb.padding()
	lb.label.SetPosition(lb.x+px, lb.y+py)
	DrawRect(lb.scene.canvas, lb.x, lb.y, lb.x+w, lb.y+h, color.Black)
	lb.label.Draw()
}

func (lb *LabelButton) Contains(x, y int) bool {
	w, h := lb.Size()
	return (lb.x <= x && x <= lb.x+w) && (lb.y <= y && y <= lb.y+h)
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
