package ui

import (
	"math"
)

type LabelButton struct {
	BasePosition
	scene  *Scene
	label  *Label
	width  float64
	hitbox Rectangle
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

func (lb *LabelButton) SetWidth(val float64) {
	lb.width = val
}

func (lb *LabelButton) Render() Rectangle {
	p := lb.padding()
	// TODO: center label horizontally
	lb.label.position = lb.position.Add(p)
	lr := lb.label.Render()
	lb.hitbox = Rectangle{
		Min: lr.Min.Sub(p),
		Max: Point{max(lr.Max.X+p.X, lb.width), lr.Max.Y + p.Y},
	}

	lb.scene.canvas.SetLineWidth(1 * lb.scene.dpi / 72.0)
	lb.scene.canvas.DrawRectangle(float64(lb.hitbox.Min.X), float64(lb.hitbox.Min.Y), float64(lb.hitbox.Max.X-lb.hitbox.Min.X), float64(lb.hitbox.Max.Y-lb.hitbox.Min.Y))
	lb.scene.canvas.Stroke()

	// drawRect(lb.scene.canvas, lb.hitbox.Min.X, lb.hitbox.Min.Y, lb.hitbox.Max.X, lb.hitbox.Max.Y, color.Black)
	return lb.hitbox
}

func (lb *LabelButton) Contains(x, y float64) bool {
	p := Point{x, y}
	return p.In(lb.hitbox)
}

// TODO: rethink the padding strategy
func (lb *LabelButton) padding() Point {
	fs := lb.label.fontSize
	// The value 72.0 corresponds to the standard reference DPI in digital typography: 1 point = 1/72 of an inch—a convention inherited from PostScript/PDF and adopted by almost all font systems (including truetype.Options.DPI, which defaults to 72 if not specified).
	// TODO: move 72.0 to another place
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
	return Point{x, y}
}
