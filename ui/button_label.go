package ui

import (
	"image"
	"image/color"
	"maps"

	"github.com/fogleman/gg"
)

var labelBtnPad = image.Point{38, 19}

type LabelButton struct {
	canvas *gg.Context
	label  *Label
	pos    image.Point
	size   image.Point
}

func NewLabelButton(c Context, text string) *LabelButton {
	lb := &LabelButton{canvas: c.canvas}
	lb.label = NewLabel(c, text)
	lb.size = lb.label.Size().Add(labelBtnPad.Mul(2))
	return lb
}

func (lb *LabelButton) Size() image.Point {
	return lb.label.Size().Add(labelBtnPad.Mul(2))
}

func (lb *LabelButton) State() map[string]any {
	dst := map[string]any{"padding": labelBtnPad}
	maps.Copy(dst, lb.label.State())
	return dst
}

func (lb *LabelButton) Render(x, y int) (dirtyRegs []image.Rectangle) {
	lb.pos = image.Pt(x, y)
	hb := lb.hitbox()
	lineWidth := 4
	lb.canvas.SetLineWidth(float64(lineWidth))
	lb.canvas.DrawRectangle(
		float64(hb.Min.X), float64(hb.Min.Y),
		float64(hb.Dx()), float64(hb.Dy()))
	lb.canvas.SetColor(color.Black)
	lb.canvas.Stroke()

	margin := image.Pt(lineWidth/2+1, lineWidth/2+1)
	dirtyReg := image.Rectangle{
		Min: hb.Min.Sub(margin),
		Max: hb.Max.Add(margin),
	}

	return append(
		lb.label.Render(x+labelBtnPad.X, y+labelBtnPad.Y),
		dirtyReg,
	)
}

func (lb *LabelButton) Contains(x, y int) bool {
	return image.Pt(x, y).In(lb.hitbox())
}

func (lb *LabelButton) hitbox() image.Rectangle {
	return image.Rectangle{
		Min: lb.pos,
		Max: lb.pos.Add(lb.size),
	}
}
