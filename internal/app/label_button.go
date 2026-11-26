package app

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/ui"
)

var labelBtnPad = image.Point{38, 19}

type labelButtonWidget struct {
	canvas *gg.Context
	label  *ui.Label
	size   image.Point
}

type LabelButton struct {
	*ui.Button
	widget *labelButtonWidget
}

func NewLabelButton(c *image.RGBA, text string, opts ...ui.LabelOption) *LabelButton {
	lb := &LabelButton{}
	lb.widget = &labelButtonWidget{canvas: gg.NewContextForRGBA(c)}
	lb.widget.label = ui.NewLabel(c, text, opts...)
	lb.widget.size = lb.widget.label.Size().Add(labelBtnPad.Mul(2))
	lb.Button = ui.NewButton(lb.widget)
	return lb
}

func (lb *labelButtonWidget) Footprint() []uint64 {
	return lb.label.Footprint()
}

func (lb *labelButtonWidget) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	hb := image.Rectangle{Min: pos, Max: pos.Add(lb.size)}
	lineWidth := 3
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

func (lb *labelButtonWidget) Size() image.Point {
	return lb.size
}
