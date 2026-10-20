package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	width, height int
	fillColor     color.Color

	canvas *gg.Context
}

func NewRectangle(c *Canvas, width, height int, fillColor color.Color) *Rectangle {
	s := &Rectangle{
		canvas:    c.context,
		width:     width,
		height:    height,
		fillColor: fillColor,
	}
	return s
}

func (r *Rectangle) Footprint() []uint64 {
	red, green, blue, alpha := r.fillColor.RGBA()
	return []uint64{uint64(r.width), uint64(r.height), uint64(red), uint64(green), uint64(blue), uint64(alpha)}
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
