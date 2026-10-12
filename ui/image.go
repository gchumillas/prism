package ui

import (
	"hash/maphash"
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

func LoadImage(c Context, r io.Reader) (i *Image, err error) {
	i = &Image{canvas: c.canvas, id: imgCounter.Add(1)}
	if i.img, _, err = image.Decode(r); err != nil {
		return
	}
	return
}

func (i *Image) Size() image.Point {
	b := i.img.Bounds()
	return image.Pt(b.Dx(), b.Dy())
}

func (i *Image) Footprint() []uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, i.id)
	return []uint64{h.Sum64()}
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
