package ui

import (
	"image"

	"github.com/fogleman/gg"
)

var paddingBtn = image.Point{38, 19}

type LabelButton struct {
	canvas *gg.Context
	hitbox image.Rectangle
	vstack *VStack
	label  *Label
}

func (s *Scene) LabelButton(text string) *LabelButton {
	lb := &LabelButton{canvas: s.canvas}
	lb.vstack = s.VStack(Padding(paddingBtn.X, paddingBtn.Y))
	lb.label = s.Label(text)
	lb.vstack.Add(lb.label)
	return lb
}

func (lb *LabelButton) X() int       { return lb.vstack.X() }
func (lb *LabelButton) SetX(val int) { lb.vstack.SetX(val) }

func (lb *LabelButton) Y() int       { return lb.vstack.Y() }
func (lb *LabelButton) SetY(val int) { lb.vstack.SetY(val) }

func (lb *LabelButton) Height() int {
	return lb.label.Height() + 2*paddingBtn.Y
}

func (lb *LabelButton) Render() image.Rectangle {
	size := lb.label.Size().Add(paddingBtn.Mul(2))
	lb.hitbox = image.Rectangle{
		Min: lb.vstack.position,
		Max: lb.vstack.position.Add(size),
	}

	lb.canvas.SetLineWidth(4)
	lb.canvas.DrawRectangle(
		float64(lb.hitbox.Min.X), float64(lb.hitbox.Min.Y),
		float64(lb.hitbox.Max.X-lb.hitbox.Min.X), float64(lb.hitbox.Max.Y-lb.hitbox.Min.Y))
	lb.canvas.Stroke()

	return lb.hitbox.Union(lb.vstack.Render())
}

func (lb *LabelButton) Contains(x, y int) bool {
	p := image.Pt(x, y)
	return p.In(lb.hitbox)
}
