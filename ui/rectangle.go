package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	Position  image.Point
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

func (r *Rectangle) SetPosition(val image.Point) {
	r.Position = val
}

func (r *Rectangle) Height() int {
	return r.size.Y
}

func (r *Rectangle) Render() []image.Rectangle {
	rec := image.Rectangle{Min: r.Position, Max: r.Position.Add(r.size)}
	r.canvas.DrawRectangle(
		float64(rec.Min.X), float64(rec.Min.Y),
		float64(rec.Dx()), float64(rec.Dy()),
	)
	r.canvas.SetColor(r.fillColor)
	r.canvas.Fill()
	return []image.Rectangle{rec}
}
