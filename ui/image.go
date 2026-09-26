package ui

import (
	"image"
	"io"

	"github.com/fogleman/gg"
)

type Image struct {
	BasePosition
	canvas *gg.Context
	img    image.Image
}

func (s *Scene) ImageFromReader(r io.Reader) (img *Image, err error) {
	img = &Image{canvas: s.canvas}
	if img.img, _, err = image.Decode(r); err != nil {
		return
	}
	return
}

func (i *Image) Render() image.Rectangle {
	i.canvas.DrawImage(i.img, i.position.X, i.position.Y)
	b := i.img.Bounds()
	size := b.Max.Sub(b.Min)

	return image.Rectangle{
		Min: i.position,
		Max: i.position.Add(size),
	}
}
