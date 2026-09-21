package ui

import (
	"image/color"
	"math"
)

type LabelButton struct {
	BasePosition
	scene *Scene
	label *Label

	width, height int
}

func (s *Scene) LabelButton(text string, opts ...LabelButtonOpt) *LabelButton {
	lb := &LabelButton{
		scene: s,
		label: s.Label(text),
	}
	for _, opt := range opts {
		opt.applyLabelButton(lb)
	}
	return lb
}

func (lb *LabelButton) Label() *Label {
	return lb.label
}

func (lb *LabelButton) Size() (int, int) {
	w, h := lb.label.Size()
	if lb.width > w {
		w = lb.width
	}
	px, py := lb.padding()
	return w + 2*px, h + 2*py
}

func (lb *LabelButton) SetSize(width, height int) {
	lb.width, lb.height = width, height
}

func (lb *LabelButton) Render() {
	w1, h1 := lb.Size()
	w0, h0 := lb.label.Size()
	lb.label.SetPosition(lb.X+(w1-w0)/2, lb.Y+(h1-h0)/2)
	drawRect(lb.scene.canvas, lb.X, lb.Y, lb.X+w1, lb.Y+h1, color.Black)
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
