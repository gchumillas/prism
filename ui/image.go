package ui

import (
	"image"
	"io"

	"github.com/fogleman/gg"
)

type Image struct {
	canvas *gg.Context
	img    image.Image
}

func LoadImage(c Context, r io.Reader) (i *Image, err error) {
	i = &Image{canvas: c.canvas}
	if i.img, _, err = image.Decode(r); err != nil {
		return
	}
	return
}

func (i *Image) Size() image.Point {
	b := i.img.Bounds()
	return image.Pt(b.Dx(), b.Dy())
}

func (i *Image) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	i.canvas.DrawImage(i.img, pos.X, pos.Y)
	b := i.img.Bounds()
	size := b.Max.Sub(b.Min)

	return []image.Rectangle{{
		Min: pos,
		Max: pos.Add(size),
	}}
}
