package ui

import (
	"image"
	"io"

	"github.com/fogleman/gg"
)

type Image struct {
	Position image.Point
	canvas   *gg.Context
	img      image.Image
}

func (c Context) LoadImage(r io.Reader) (i *Image, err error) {
	i = &Image{canvas: c.canvas}
	if i.img, _, err = image.Decode(r); err != nil {
		return
	}
	return
}

// TODO: replace with SetPosition
func (i *Image) SetX(val int) { i.Position.X = val }
func (i *Image) SetY(val int) { i.Position.Y = val }

func (i *Image) Size() image.Point {
	b := i.img.Bounds()
	return image.Pt(b.Dx(), b.Dy())
}

func (i *Image) Height() int {
	return i.Size().Y
}

func (i *Image) Render() []image.Rectangle {
	i.canvas.DrawImage(i.img, i.Position.X, i.Position.Y)
	b := i.img.Bounds()
	size := b.Max.Sub(b.Min)

	return []image.Rectangle{{
		Min: i.Position,
		Max: i.Position.Add(size),
	}}
}
