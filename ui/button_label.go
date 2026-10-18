package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

var labelBtnPad = image.Point{38, 19}

type labelButton struct {
	canvas *gg.Context
	label  *Label
	size   image.Point
}

func newLabelButton(c *Canvas, text string, opts ...labelOption) *labelButton {
	lb := &labelButton{canvas: c.context}
	lb.label = NewLabel(c, text, opts...)
	lb.size = lb.label.size().Add(labelBtnPad.Mul(2))
	return lb
}

func AttachLabelButtonTo(parent *Stack, text string, opts ...labelOption) *Button {
	return AttachButtonTo(parent, newLabelButton(parent.canvas, text, opts...))
}

func (lb *labelButton) Footprint() []uint64 {
	return lb.label.Footprint()
}

func (lb *labelButton) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	hb := image.Rectangle{Min: pos, Max: pos.Add(lb.size)}
	lineWidth := 4
	lb.canvas.SetLineWidth(float64(lineWidth))
	lb.canvas.DrawRectangle(
		float64(hb.Min.X), float64(hb.Min.Y),
		float64(hb.Dx()), float64(hb.Dy()))
	lb.canvas.SetColor(color.Black)
	lb.canvas.Stroke()

	// addSafeMargin(r, margin)
	margin := image.Pt(lineWidth/2+1, lineWidth/2+1)
	// TODO: use an Outset(rect, margin) utility instead (affect others)
	dirtyReg := image.Rectangle{
		Min: hb.Min.Sub(margin),
		Max: hb.Max.Add(margin),
	}

	dr, _ := lb.label.Render(x+labelBtnPad.X, y+labelBtnPad.Y)
	return append(dr, dirtyReg), image.Pt(hb.Dx(), hb.Dy())
}
