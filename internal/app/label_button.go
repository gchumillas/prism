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

func (lbw *labelButtonWidget) Footprint() []uint64 {
	return lbw.label.Footprint()
}

func (lbw *labelButtonWidget) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	hb := image.Rectangle{Min: pos, Max: pos.Add(lbw.size)}
	lineWidth := 3
	lbw.canvas.SetLineWidth(float64(lineWidth))
	lbw.canvas.DrawRectangle(
		float64(hb.Min.X), float64(hb.Min.Y),
		float64(hb.Dx()), float64(hb.Dy()))
	lbw.canvas.SetColor(color.White)
	lbw.canvas.FillPreserve()
	lbw.canvas.SetColor(color.Black)
	lbw.canvas.Stroke()

	dirtyReg := image.Rectangle{Min: hb.Min, Max: hb.Max}.Inset(-lineWidth/2 - 1)
	dr := lbw.label.Render(x+labelBtnPad.X, y+labelBtnPad.Y)
	return append(dr, dirtyReg)
}

func (lbw *labelButtonWidget) Size() image.Point {
	return lbw.size
}
