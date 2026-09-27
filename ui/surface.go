package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

type Surface struct {
	BasePosition
	canvas    *gg.Context
	size      image.Point
	fillColor color.Color
}

type SurfaceOpt interface {
	applySurface(*Surface)
}

// avail options
func (val xOpt) applySurface(x *Surface) { x.SetX(int(val)) }
func (val yOpt) applySurface(x *Surface) { x.SetY(int(val)) }

func (c Context) Surface(width, height int, fillColor color.Color, opts ...SurfaceOpt) *Surface {
	s := &Surface{
		canvas:    c.canvas,
		size:      image.Pt(width, height),
		fillColor: fillColor,
	}
	for _, opt := range opts {
		opt.applySurface(s)
	}
	return s
}

func (s *Surface) Height() int {
	return s.size.Y
}

func (s *Surface) Render() image.Rectangle {
	r := image.Rectangle{Min: s.position, Max: s.position.Add(s.size)}
	s.canvas.DrawRectangle(
		float64(r.Min.X), float64(r.Min.Y),
		float64(r.Dx()), float64(r.Dy()),
	)
	s.canvas.SetColor(s.fillColor)
	s.canvas.Fill()
	return r
}
