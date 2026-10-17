package ui

import (
	"hash/maphash"
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	BaseWidget
	width, height int
	fillColor     color.Color

	canvas *gg.Context
}

func NewRectangle(c Context, width, height int, fillColor color.Color) *Rectangle {
	s := &Rectangle{
		canvas:    c.canvas,
		width:     width,
		height:    height,
		fillColor: fillColor,
	}
	return s
}

func (r *Rectangle) Footprint() []uint64 {
	red, green, blue, alpha := r.fillColor.RGBA()

	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, r.width)
	maphash.WriteComparable(&h, r.height)
	maphash.WriteComparable(&h, red)
	maphash.WriteComparable(&h, green)
	maphash.WriteComparable(&h, blue)
	maphash.WriteComparable(&h, alpha)
	return []uint64{h.Sum64()}
}

func (r *Rectangle) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	size = image.Pt(r.width, r.height)
	rec := image.Rectangle{Min: pos, Max: pos.Add(size)}
	r.canvas.DrawRectangle(
		float64(rec.Min.X), float64(rec.Min.Y),
		float64(rec.Dx()), float64(rec.Dy()),
	)
	r.canvas.SetColor(r.fillColor)
	r.canvas.Fill()
	return []image.Rectangle{rec}, size
}
