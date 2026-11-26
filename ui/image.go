package ui

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"sync/atomic"
)

var imgCounter atomic.Uint64

type Image struct {
	canvas draw.Image
	img    image.Image
	id     uint64
}

func LoadImageFromPNG(canvas draw.Image, r io.Reader) (*Image, error) {
	var err error
	img := &Image{canvas: canvas, id: imgCounter.Add(1)}
	if img.img, err = png.Decode(r); err != nil {
		return nil, err
	}
	return img, err
}

func LoadImageFromJPEG(canvas draw.Image, r io.Reader) (*Image, error) {
	var err error
	img := &Image{canvas: canvas, id: imgCounter.Add(1)}
	if img.img, err = jpeg.Decode(r); err != nil {
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

func (i *Image) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	size := i.Size()
	draw.Draw(i.canvas, image.Rectangle{Min: pos, Max: pos.Add(size)}, i.img, i.img.Bounds().Min, draw.Over)
	// TODO: add safe margin (affect others)
	return []image.Rectangle{{Min: pos, Max: pos.Add(size)}}
}
