package ui_test

import (
	"hash/maphash"
	"image"
	"image/color"
	"image/draw"

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

// ---- Rect ----

type Rect struct {
	canvas draw.Image
	image.Rectangle
	color color.Color
}

func NewRect(canvas draw.Image, x0, y0, x1, y1 int, c color.Color) *Rect {
	return &Rect{
		canvas:    canvas,
		Rectangle: image.Rect(x0, y0, x1, y1),
		color:     c,
	}
}

func (r *Rect) Footprint() []uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, r.Dx())
	maphash.WriteComparable(&h, r.Dy())
	maphash.WriteComparable(&h, r.color)
	return []uint64{h.Sum64()}
}

func (r *Rect) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	size = r.Size()
	r.Rectangle = image.Rectangle{Min: pos, Max: pos.Add(size)}
	draw.Draw(r.canvas, r.Rectangle, image.NewUniform(r.color), image.Pt(0, 0), draw.Src)
	return []image.Rectangle{r.Rectangle}, size
}

type Rect1 struct {
	*Rect
}

func NewRect1(canvas draw.Image, x0, y0, x1, y1 int, c color.Color) *Rect1 {
	return &Rect1{Rect: NewRect(canvas, x0, y0, x1, y1, c)}
}
