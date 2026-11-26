package app

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	canvas      *gg.Context
	width       int
	height      int
	lineWidth   float64
	fillColor   color.Color
	strokeColor color.Color
}

type RectangleOption interface {
	ApplyRectangle(*Rectangle)
}

func (val WithFillColor) ApplyRectangle(r *Rectangle)   { r.fillColor = uint32ToColor(uint32(val)) }
func (val WithStrokeColor) ApplyRectangle(r *Rectangle) { r.strokeColor = uint32ToColor(uint32(val)) }
func (val WithLineWidth) ApplyRectangle(r *Rectangle)   { r.lineWidth = float64(val) }

func NewRectangle(c *image.RGBA, width, height int, opts ...RectangleOption) *Rectangle {
	s := &Rectangle{
		canvas:      gg.NewContextForRGBA(c),
		width:       width,
		height:      height,
		lineWidth:   3,
		fillColor:   color.Transparent,
		strokeColor: color.Transparent,
	}
	for _, opt := range opts {
		opt.ApplyRectangle(s)
	}
	return s
}

func (r *Rectangle) Footprint() []uint64 {
	r1, g1, b1, a1 := r.fillColor.RGBA()
	r2, g2, b2, a2 := r.strokeColor.RGBA()
	return []uint64{
		uint64(r.width), uint64(r.height), math.Float64bits(r.lineWidth),
		uint64(r1), uint64(g1), uint64(b1), uint64(a1),
		uint64(r2), uint64(g2), uint64(b2), uint64(a2),
	}
}

func (r *Rectangle) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	rec := image.Rectangle{Min: pos, Max: pos.Add(r.Size())}
	r.canvas.DrawRectangle(
		float64(rec.Min.X), float64(rec.Min.Y),
		float64(rec.Dx()), float64(rec.Dy()),
	)
	r.canvas.SetColor(r.fillColor)
	r.canvas.FillPreserve()
	r.canvas.SetLineWidth(r.lineWidth)
	r.canvas.SetColor(r.strokeColor)
	r.canvas.Stroke()
	pad := int(math.Ceil(r.lineWidth/2 + 1))
	return []image.Rectangle{rec.Inset(-pad)}
}

func (r *Rectangle) Size() image.Point {
	return image.Pt(r.width, r.height)
}
