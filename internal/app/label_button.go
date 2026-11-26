package app

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/ui"
)

var labelBtnPad = image.Point{38, 19}

type LabelButton struct {
	canvas *gg.Context
	label  *ui.Label
	size   image.Point
}

func NewLabelButton(c *ui.Canvas, text string, opts ...ui.LabelOption) *ui.Button {
	lb := &LabelButton{canvas: gg.NewContextForRGBA(c.Image())}
	lb.label = ui.NewLabel(c, text, opts...)
	lb.size = lb.label.Size().Add(labelBtnPad.Mul(2))
	return ui.NewButton(lb)
}

func (lb *LabelButton) Footprint() []uint64 {
	return lb.label.Footprint()
}

func (lb *LabelButton) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	hb := image.Rectangle{Min: pos, Max: pos.Add(lb.size)}
	lineWidth := 4
	lb.canvas.SetLineWidth(float64(lineWidth))
	lb.canvas.DrawRectangle(
		float64(hb.Min.X), float64(hb.Min.Y),
		float64(hb.Dx()), float64(hb.Dy()))
	lb.canvas.SetColor(color.White)
	lb.canvas.FillPreserve()
	lb.canvas.SetColor(color.Black)
	lb.canvas.Stroke()

	dirtyReg := image.Rectangle{Min: hb.Min, Max: hb.Max}.Inset(-lineWidth/2 - 1)
	dr := lb.label.Render(x+labelBtnPad.X, y+labelBtnPad.Y)
	return append(dr, dirtyReg)
}

func (lb *LabelButton) Size() image.Point {
	return lb.size
}
