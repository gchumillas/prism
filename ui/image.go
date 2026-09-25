package ui

import (
	"image"
	"io"
)

type Image struct {
	BasePosition
	scene *Scene
	img   image.Image
}

func (s *Scene) ImageFromReader(r io.Reader) (img *Image, err error) {
	img = &Image{scene: s}
	if img.img, _, err = image.Decode(r); err != nil {
		return
	}
	return
}

func (i *Image) Render() image.Rectangle {
	i.scene.canvas.DrawImage(i.img, i.position.X, i.position.Y)
	b := i.img.Bounds()
	size := b.Max.Sub(b.Min)

	return image.Rectangle{
		Min: i.position,
		Max: i.position.Add(size),
	}
}
