package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	BasePosition
	canvas    *gg.Context
	size      image.Point
	fillColor color.Color
}

type RectangleOpt interface {
	applyRectangle(*Rectangle)
}

// avail options
func (val xOpt) applyRectangle(x *Rectangle) { x.SetX(int(val)) }
func (val yOpt) applyRectangle(x *Rectangle) { x.SetY(int(val)) }

func (c Context) Rectangle(width, height int, fillColor color.Color, opts ...RectangleOpt) *Rectangle {
	s := &Rectangle{
		canvas:    c.canvas,
		size:      image.Pt(width, height),
		fillColor: fillColor,
	}
	for _, opt := range opts {
		opt.applyRectangle(s)
	}
	return s
}

func (r *Rectangle) Height() int {
	return r.size.Y
}

func (r *Rectangle) Render() []image.Rectangle {
	rec := image.Rectangle{Min: r.position, Max: r.position.Add(r.size)}
	r.canvas.DrawRectangle(
		float64(rec.Min.X), float64(rec.Min.Y),
		float64(rec.Dx()), float64(rec.Dy()),
	)
	r.canvas.SetColor(r.fillColor)
	r.canvas.Fill()
	return []image.Rectangle{rec}
}
