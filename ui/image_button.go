package ui

import (
	"hash/maphash"
	"image"
	"image/color"
	"io"

	"github.com/fogleman/gg"
)

var imgBtnPad = image.Pt(5, 5)

type ImageButton struct {
	canvas *gg.Context
	img    *Image
	pos    image.Point
	size   image.Point
}

func LoadImageButton(c Context, r io.Reader) (_ *ImageButton, err error) {
	ib := &ImageButton{canvas: c.canvas}
	if ib.img, err = LoadImage(c, r); err != nil {
		return
	}
	ib.size = ib.img.size().Add(imgBtnPad.Mul(2))
	return ib, nil
}

func (ib *ImageButton) Footprint() []uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, ib.pos)
	return append([]uint64{h.Sum64()}, ib.img.Footprint()...)
}

func (ib *ImageButton) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	ib.pos = image.Pt(x, y)
	hb := ib.hitbox()
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

func (ib *ImageButton) Contains(x, y int) bool {
	return image.Pt(x, y).In(ib.hitbox())
}

func (ib *ImageButton) hitbox() image.Rectangle {
	return image.Rectangle{
		Min: ib.pos,
		Max: ib.pos.Add(ib.size),
	}
}
