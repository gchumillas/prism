package ui

import (
	"image"
	"io"
	"sync/atomic"

	"github.com/fogleman/gg"
)

var imgCounter atomic.Uint64

type Image struct {
	canvas *gg.Context
	img    image.Image
	id     uint64
}

func LoadImage(c *Canvas, r io.Reader) (*Image, error) {
	var err error
	img := &Image{canvas: c.context, id: imgCounter.Add(1)}
	if img.img, _, err = image.Decode(r); err != nil {
		return nil, err
	}
	return img, err
}

func (i *Image) Size() image.Point {
	b := i.img.Bounds()
	return image.Pt(b.Dx(), b.Dy())
}

func (i *Image) Footprint() []uint64 {
	return []uint64{i.id}
}

func (i *Image) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	i.canvas.DrawImage(i.img, pos.X, pos.Y)
	b := i.img.Bounds()
	size = image.Pt(b.Dx(), b.Dy())
	// TODO: add safe margin (affect others)
	return []image.Rectangle{{Min: pos, Max: pos.Add(size)}}, size
}
