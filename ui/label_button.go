package ui

import (
	"image"
	"image/color"
	"math"
)

type LabelButton struct {
	BasePosition
	scene *Scene
	label *Label
	width int

	rect image.Rectangle
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

func (lb *LabelButton) SetWidth(val int) {
	lb.width = val
}

func (lb *LabelButton) Render() image.Rectangle {
	p := lb.padding()
	lb.label.x, lb.label.y = lb.x+p.X, lb.y+p.Y
	lr := lb.label.Render()
	lb.rect = image.Rectangle{
		Min: lr.Min.Sub(p),
		Max: image.Pt(max(lr.Max.X+p.X, lb.width), lr.Max.Y+p.Y),
	}
	drawRect(lb.scene.canvas, lb.rect.Min.X, lb.rect.Min.Y, lb.rect.Max.X, lb.rect.Max.Y, color.Black)
	return lb.rect
}

func (lb *LabelButton) Contains(x, y int) bool {
	p := image.Pt(x, y)
	return p.In(lb.rect)
}

func (lb *LabelButton) padding() image.Point {
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
	return image.Pt(int(x), int(y))
}
