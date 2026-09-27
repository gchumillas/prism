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

type ImageOpt interface {
	applyImage(*Image)
}

// avail options
func (val xOpt) applyImage(x *Image) { x.SetX(int(val)) }
func (val yOpt) applyImage(x *Image) { x.SetY(int(val)) }

func (c Context) ImageFromReader(r io.Reader, opts ...ImageOpt) (i *Image, err error) {
	i = &Image{canvas: c.canvas}
	for _, opt := range opts {
		opt.applyImage(i)
	}
	if i.img, _, err = image.Decode(r); err != nil {
		return
	}
	return
}

func (i *Image) Size() image.Point {
	b := i.img.Bounds()
	return image.Pt(b.Dx(), b.Dy())
}

func (i *Image) Height() int {
	return i.Size().Y
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
