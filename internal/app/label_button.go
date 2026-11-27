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
	Width  int
}

func NewLabelButton(c *image.RGBA, text string, opts ...ui.LabelOption) *ui.Button[*labelButtonWidget] {
	lbw := &labelButtonWidget{canvas: gg.NewContextForRGBA(c)}
	lbw.label = ui.NewLabel(c, text, opts...)
	lbw.size = lbw.label.Size().Add(labelBtnPad.Mul(2))
	return ui.NewButton(lbw)
}

func (lbw *labelButtonWidget) Footprint() []uint64 {
	return append([]uint64{uint64(lbw.Width)}, lbw.label.Footprint()...)
}

func (lbw *labelButtonWidget) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	hb := image.Rectangle{Min: pos, Max: pos.Add(lbw.Size())}
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
	size := lbw.size
	size.X = max(lbw.Width, lbw.size.X)
	return size
}
