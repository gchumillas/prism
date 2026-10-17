package ui

import (
	"image"
	"image/color"
	"io"

	"github.com/fogleman/gg"
)

var imgBtnPad = image.Pt(5, 5)

type imageButton struct {
	canvas *gg.Context
	img    *Image
	size   image.Point
}

func loadImageButton(c Context, r io.Reader) (_ *imageButton, err error) {
	ib := &imageButton{canvas: c.canvas}
	if ib.img, err = LoadImage(c, r); err != nil {
		return
	}
	ib.size = ib.img.size().Add(imgBtnPad.Mul(2))
	return ib, nil
}

func AttachImageButton(owner Stack, ctx Context, r io.Reader) (*Button, error) {
	img, err := loadImageButton(ctx, r)
	if err != nil {
		return nil, err
	}
	btn := newClickableWidget(img)
	owner.Add(btn)
	return btn.Button, nil
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

	margin := image.Pt(lineWidth/2+1, lineWidth/2+1)
	dirtyReg := image.Rectangle{
		Min: hb.Min.Sub(margin),
		Max: hb.Max.Add(margin),
	}

	dr, _ := ib.img.Render(x+imgBtnPad.X, y+imgBtnPad.Y)
	return append(dr, dirtyReg), image.Pt(hb.Dx(), hb.Dy())
}
