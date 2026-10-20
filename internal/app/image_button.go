package app

import (
	"image"
	"image/color"
	"io"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/ui"
)

var imgBtnPad = image.Pt(5, 5)

type imageButton struct {
	canvas *gg.Context
	img    *ui.Image
	size   image.Point
}

func LoadImageButton(c *ui.Canvas, r io.Reader) (_ *ui.Button, err error) {
	ib := &imageButton{canvas: gg.NewContextForRGBA(c.Image())}
	if ib.img, err = ui.LoadImage(c, r); err != nil {
		return
	}
	ib.size = ib.img.Size().Add(imgBtnPad.Mul(2))
	return ui.NewButton(ib), nil
}

func (ib *imageButton) Footprint() []uint64 {
	return ib.img.Footprint()
}

func (ib *imageButton) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	hb := image.Rectangle{Min: pos, Max: pos.Add(ib.size)}
	lineWidth := 4
	ib.canvas.SetLineWidth(float64(lineWidth))
	ib.canvas.SetColor(color.Black)
	ib.canvas.DrawRectangle(
		float64(hb.Min.X), float64(hb.Min.Y),
		float64(hb.Dx()), float64(hb.Dy()))
	ib.canvas.Stroke()

	dirtyReg := image.Rectangle{Min: hb.Min, Max: hb.Max}.Inset(-lineWidth/2 - 1)
	dr, _ := ib.img.Render(x+imgBtnPad.X, y+imgBtnPad.Y)
	return append(dr, dirtyReg), image.Pt(hb.Dx(), hb.Dy())
}
