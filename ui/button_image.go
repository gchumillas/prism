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
	vstack  *VStack
	img     *Image
}

func LoadImageButton(c Context, r io.Reader) (_ *ImageButton, err error) {
	ib := &ImageButton{canvas: c.canvas}
	ib.vstack = NewVStack(c)
	ib.Padding = imgBtnPad
	if ib.img, err = LoadImage(c, r); err != nil {
		return
	}
	ib.vstack.Add(ib.img)
	return ib, nil
}

func (ib *ImageButton) Size() image.Point {
	return ib.img.Size().Add(imgBtnPad.Mul(2))
}

func (ib *ImageButton) Render(pos image.Point) (dirtyRegs []image.Rectangle) {
	ib.vstack.Padding = ib.Padding
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
	return append(ib.vstack.Render(pos), image.Rectangle{
		Min: ib.hitbox.Min.Sub(margin),
		Max: ib.hitbox.Max.Add(margin),
	})
}

func (ib *ImageButton) Contains(p image.Point) bool {
	if ib.hitbox.Empty() {
		panic("The hitbox is empty. Did you forget calling Render?")
	}
	return p.In(ib.hitbox)
}
