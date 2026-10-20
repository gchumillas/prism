package app

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/ui"
)

var labelBtnPad = image.Point{38, 19}

type labelButton struct {
	canvas *gg.Context
	label  *ui.Label
	size   image.Point
}

func NewLabelButton(c *ui.Canvas, text string, opts ...ui.LabelOption) *ui.Button {
	lb := &labelButton{canvas: gg.NewContextForRGBA(c.Image())}
	lb.label = ui.NewLabel(c, text, opts...)
	lb.size = lb.label.Size().Add(labelBtnPad.Mul(2))
	return ui.NewButton(lb)
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

	dirtyReg := image.Rectangle{Min: hb.Min, Max: hb.Max}.Inset(-lineWidth/2 - 1)
	dr, _ := lb.label.Render(x+labelBtnPad.X, y+labelBtnPad.Y)
	return append(dr, dirtyReg), image.Pt(hb.Dx(), hb.Dy())
}
