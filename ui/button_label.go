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
	label   *Label
}

func NewLabelButton(c Context, text string) *LabelButton {
	lb := &LabelButton{canvas: c.canvas}
	lb.Padding = labelBtnPad
	lb.label = NewLabel(c, text)
	return lb
}

func (lb *LabelButton) Size() image.Point {
	return lb.label.Size().Add(labelBtnPad.Mul(2))
}

func (lb *LabelButton) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
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
	dirtyReg := image.Rectangle{
		Min: lb.hitbox.Min.Sub(margin),
		Max: lb.hitbox.Max.Add(margin),
	}

	return append(
		lb.label.Render(x+lb.Padding.X, y+lb.Padding.Y),
		dirtyReg,
	)
}

func (lb *LabelButton) Contains(x, y int) bool {
	if lb.hitbox.Empty() {
		panic(ErrEmptyHitbox)
	}
	return image.Pt(x, y).In(lb.hitbox)
}
