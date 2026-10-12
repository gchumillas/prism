package ui

import (
	"hash/maphash"
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	Width, Height int
	FillColor     color.Color

	canvas *gg.Context
}

func NewRectangle(c Context, width, height int, fillColor color.Color) *Rectangle {
	s := &Rectangle{
		Width:     width,
		Height:    height,
		canvas:    c.canvas,
		FillColor: fillColor,
	}
	return s
}

func (r *Rectangle) Size() image.Point {
	return image.Pt(r.Width, r.Height)
}

func (r *Rectangle) Footprint() []uint64 {
	red, green, blue, alpha := r.FillColor.RGBA()

	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, r.Width)
	maphash.WriteComparable(&h, r.Height)
	maphash.WriteComparable(&h, red)
	maphash.WriteComparable(&h, green)
	maphash.WriteComparable(&h, blue)
	maphash.WriteComparable(&h, alpha)
	return []uint64{h.Sum64()}
}

func (r *Rectangle) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	size := image.Pt(r.Width, r.Height)
	rec := image.Rectangle{Min: pos, Max: pos.Add(size)}
	r.canvas.DrawRectangle(
		float64(rec.Min.X), float64(rec.Min.Y),
		float64(rec.Dx()), float64(rec.Dy()),
	)
	r.canvas.SetColor(r.FillColor)
	r.canvas.Fill()
	return []image.Rectangle{rec}
}
