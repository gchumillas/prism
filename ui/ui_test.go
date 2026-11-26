package ui_test

import (
	"hash/maphash"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/ui"

	_ "image/png"
)

// ---- common variables and constants ----

var seed = maphash.MakeSeed()

// ---- Env ----

type Env struct {
	image  *image.RGBA
	canvas *ui.Canvas
}

func NewEnv() *Env {
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	c := ui.NewCanvasFromRGBA(img)
	return &Env{image: img, canvas: c}
}

// ---- Rectangle ----

type Rectangle struct {
	width, height int
	fillColor     color.Color

	canvas *gg.Context
}

func NewRectangle(c *ui.Canvas, width, height int, fillColor color.Color) *Rectangle {
	s := &Rectangle{
		canvas:    gg.NewContextForRGBA(c.Image()),
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

func (r *Rectangle) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	rec := image.Rectangle{Min: pos, Max: pos.Add(r.Size())}
	r.canvas.DrawRectangle(
		float64(rec.Min.X), float64(rec.Min.Y),
		float64(rec.Dx()), float64(rec.Dy()),
	)
	r.canvas.SetColor(r.fillColor)
	r.canvas.Fill()
	return []image.Rectangle{rec}
}

func (r *Rectangle) Size() image.Point {
	return image.Pt(r.width, r.height)
}
