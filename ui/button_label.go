package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

var labelBtnPad = image.Point{38, 19}

type LabelButton struct {
	Padding image.Point
	canvas  *gg.Context
	hitbox  image.Rectangle
	vstack  *VStack
	label   *Label
}

func NewLabelButton(c Context, text string) *LabelButton {
	lb := &LabelButton{canvas: c.canvas}
	lb.vstack = NewVStack(c)
	lb.Padding = labelBtnPad
	lb.label = NewLabel(c, text)
	lb.vstack.Add(lb.label)
	return lb
}

func (lb *LabelButton) Size() image.Point {
	return lb.label.Size().Add(labelBtnPad.Mul(2))
}

func (lb *LabelButton) Render(pos image.Point) (dirtyRegs []image.Rectangle) {
	lb.vstack.Padding = lb.Padding
	size := lb.label.Size().Add(labelBtnPad.Mul(2))
	lb.hitbox = image.Rectangle{
		Min: pos,
		Max: pos.Add(size),
	}

	lineWidth := 4
	lb.canvas.SetLineWidth(float64(lineWidth))
	lb.canvas.DrawRectangle(
		float64(lb.hitbox.Min.X), float64(lb.hitbox.Min.Y),
		float64(lb.hitbox.Dx()), float64(lb.hitbox.Dy()))
	lb.canvas.SetColor(color.Black)
	lb.canvas.Stroke()

	margin := image.Pt(lineWidth/2+1, lineWidth/2+1)
	return append(lb.vstack.Render(pos), image.Rectangle{
		Min: lb.hitbox.Min.Sub(margin),
		Max: lb.hitbox.Max.Add(margin),
	})
}

func (lb *LabelButton) Contains(p image.Point) bool {
	if lb.hitbox.Empty() {
		panic(ErrEmptyHitbox)
	}
	return p.In(lb.hitbox)
}
