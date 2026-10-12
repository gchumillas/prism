package ui

import (
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

func (r *Rectangle) Footprint() []any {
	red, green, blue, alpha := r.FillColor.RGBA()
	return []any{r.Width, r.Height, [4]uint32{red, green, blue, alpha}}
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
