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

func (i *Image) Render() Rectangle {
	i.scene.canvas.DrawImage(i.img, int(i.position.X), int(i.position.Y))
	b := i.img.Bounds()
	size := b.Max.Sub(b.Min)

	return Rectangle{
		Min: i.position,
		Max: i.position.Add(Point{X: float64(size.X), Y: float64(size.Y)}),
	}
}
