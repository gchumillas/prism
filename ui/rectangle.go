package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	canvas    *gg.Context
	size      image.Point
	fillColor color.Color
}

func NewRectangle(c Context, width, height int, fillColor color.Color) *Rectangle {
	s := &Rectangle{
		canvas:    c.canvas,
		size:      image.Pt(width, height),
		fillColor: fillColor,
	}
	return s
}

func (r *Rectangle) Size() image.Point {
	return r.size
}

func (r *Rectangle) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	rec := image.Rectangle{Min: pos, Max: pos.Add(r.size)}
	r.canvas.DrawRectangle(
		float64(rec.Min.X), float64(rec.Min.Y),
		float64(rec.Dx()), float64(rec.Dy()),
	)
	r.canvas.SetColor(r.fillColor)
	r.canvas.Fill()
	return []image.Rectangle{rec}
}
