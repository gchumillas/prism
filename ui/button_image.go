package ui

import (
	"image"
	"image/color"
	"io"

	"github.com/fogleman/gg"
)

var imgBtnPad = image.Pt(5, 5)

type ImageButton struct {
	Padding image.Point
	canvas  *gg.Context
	hitbox  image.Rectangle
	img     *Image
}

func LoadImageButton(c Context, r io.Reader) (_ *ImageButton, err error) {
	ib := &ImageButton{canvas: c.canvas}
	ib.Padding = imgBtnPad
	if ib.img, err = LoadImage(c, r); err != nil {
		return
	}
	return ib, nil
}

func (ib *ImageButton) Size() image.Point {
	return ib.img.Size().Add(imgBtnPad.Mul(2))
}

func (ib *ImageButton) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	size := ib.img.Size().Add(imgBtnPad.Mul(2))
	ib.hitbox = image.Rectangle{
		Min: pos,
		Max: pos.Add(size),
	}

	lineWidth := 4
	ib.canvas.SetLineWidth(float64(lineWidth))
	ib.canvas.SetColor(color.Black)
	ib.canvas.DrawRectangle(
		float64(ib.hitbox.Min.X), float64(ib.hitbox.Min.Y),
		float64(ib.hitbox.Dx()), float64(ib.hitbox.Dy()))
	ib.canvas.Stroke()

	margin := image.Pt(lineWidth/2+1, lineWidth/2+1)
	dirtyReg := image.Rectangle{
		Min: ib.hitbox.Min.Sub(margin),
		Max: ib.hitbox.Max.Add(margin),
	}

	return append(
		ib.img.Render(x+ib.Padding.X, y+ib.Padding.Y),
		dirtyReg,
	)
}

func (ib *ImageButton) Contains(x, y int) bool {
	if ib.hitbox.Empty() {
		panic(ErrEmptyHitbox)
	}
	return image.Pt(x, y).In(ib.hitbox)
}
