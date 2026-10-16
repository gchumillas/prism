//go:build regression

package ui_test

import (
	"hash/maphash"
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/internal/test"
	"github.com/gchumillas/prism/ui"
)

type Square struct {
	ctx   *gg.Context
	len   int
	bSize int
}

func NewSquare(canvas *image.RGBA, len, bSize int) *Square {
	ctx := gg.NewContextForRGBA(canvas)
	return &Square{ctx: ctx, len: len, bSize: bSize}
}

func (sq *Square) Footprint() []uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, sq.len)
	maphash.WriteComparable(&h, sq.bSize)
	return []uint64{h.Sum64()}
}

func (sq *Square) Render(x, y int) ([]image.Rectangle, image.Point) {
	sq.ctx.SetColor(color.White)
	sq.ctx.SetLineWidth(float64(sq.bSize))
	sq.ctx.DrawLine(float64(x), float64(y), float64(x+sq.len), float64(y))
	sq.ctx.DrawLine(float64(x), float64(y+sq.len), float64(x+sq.len), float64(y+sq.len))
	sq.ctx.DrawLine(float64(x), float64(y), float64(x), float64(y+sq.len))
	sq.ctx.DrawLine(float64(x+sq.len), float64(y), float64(x+sq.len), float64(y+sq.len))
	sq.ctx.Stroke()

	safeMargin := sq.bSize/2 + 1
	size := image.Pt(sq.len, sq.len)
	dirtyReg := image.Rectangle{
		Min: image.Pt(x-safeMargin, y-safeMargin),
		Max: image.Pt(x+sq.len+safeMargin, y+sq.len+safeMargin),
	}

	return []image.Rectangle{dirtyReg}, size
}

func Test(t *testing.T) {
	env := NewEnv()
	sq1 := NewSquare(env.canvas, 50, 10) // square with a border
	sq2 := NewSquare(env.canvas, 50, 0)  // square without a border
	hs := ui.NewHStack(env.ctx, ui.WithGap(50))
	hs.Add(sq1)
	hs.Add(sq2)
	hs.Render(10, 10)

	// clear the squares and repaint
	hs.Remove(sq1)
	hs.Remove(sq2)
	dirtyRegs, _ := hs.Render(10, 10)

	// clear "dirty regions"
	for _, reg := range dirtyRegs {
		draw.Draw(env.canvas, reg, image.NewUniform(color.Transparent), image.Pt(0, 0), draw.Src)
	}
	test.TakeSnapshot(t, env.canvas)

	// the canvas should be completely clean
	b := env.canvas.Bounds()
	test.RectIsColor(t, env.canvas, b.Min.X, b.Min.Y, b.Max.X, b.Max.Y, color.Transparent)
}
