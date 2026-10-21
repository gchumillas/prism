//go:build regression

package ui

import (
	"hash/maphash"
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/internal/test"
	"github.com/stretchr/testify/require"
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

// https://github.com/gchumillas/prism/pull/77
func TestPull77(t *testing.T) {
	// https://github.com/gchumillas/prism/pull/77#discussion_r4302316204
	t.Run("footprint must include direction", func(t *testing.T) {
		c := NewCanvas(320, 240)
		hs := NewVStack(c)
		vs := NewHStack(c)
		require.NotEqual(t, hs.footprint(), vs.footprint())
	})
}

func Test1(t *testing.T) {
	c := NewCanvas(320, 250)

	hs := NewHStack(c)
	hs1 := NewHStack(c)
	hs1.Add(NewRectangle(c, 50, 50, color.White))

	hs1.AttachTo(hs)
	hs.Render(0, 0)

	hs1.DetachFrom(hs)
	hs.Render(0, 0)

	hs1.AttachTo(hs)
	regs, _ := hs.Render(0, 0)
	require.NotEmpty(t, regs)
}

// https://github.com/gchumillas/prism/issues/110
func TestIssue110(t *testing.T) {
	width := 50
	x0, y0 := 100, 50
	x1, y1 := x0+50, y0+50
	c := NewCanvas(320, 240)

	// internal stack
	hs1 := NewHStack(c)
	r := NewRectangle(c, width, width, color.White)
	hs1.Add(r)

	// render a external stack that contains a internal stack
	hs := NewHStack(c)
	_ = hs1.AttachTo(hs)
	hs.Render(x0, y0)

	// remove the internal stack, render and clear "dirty regions"
	_ = hs1.DetachFrom(hs)
	dirtyRegs, _ := hs.Render(x0, y0)
	for _, reg := range dirtyRegs {
		draw.Draw(c.image, reg, image.NewUniform(color.Transparent), image.Point{}, draw.Src)
	}
	test.RectIsColor(t, c.image, x0, y0, x1, y1, color.Transparent) // has been cleared correctly

	// add the internal stack again and render
	_ = hs1.AttachTo(hs)
	hs.Render(x0, y0)
	test.RectIsColor(t, c.image, x0+10, y0+10, x1-10, y1-10, color.White) // has been rendered correctly

	// remove the internal stack again, repaint and clear "dirty regions"
	hs1.DetachFrom(hs)
	dirtyRegs, _ = hs.Render(x0, y0)
	for _, reg := range dirtyRegs {
		draw.Draw(c.image, reg, image.NewUniform(color.Transparent), image.Point{}, draw.Src)
	}
	test.RectIsColor(t, c.image, x0, y0, x1, y1, color.Transparent) // has been cleared correctly
}

// https://github.com/gchumillas/prism/issues/73
func TestIssue73(t *testing.T) {
	t.Run("dirty regions cover all pixels", func(t *testing.T) {
		c := NewCanvas(320, 240)
		sq1 := NewSquare(c.image, 50, 10) // square with a border
		sq2 := NewSquare(c.image, 50, 0)  // square without a border
		hs := NewHStack(c, WithGap(50))
		hs.Add(sq1)
		hs.Add(sq2)
		hs.Render(10, 10)

		// clear the squares and repaint
		hs.Remove(sq1)
		hs.Remove(sq2)
		dirtyRegs, _ := hs.Render(10, 10)

		// clear "dirty regions"
		for _, reg := range dirtyRegs {
			draw.Draw(c.image, reg, image.NewUniform(color.Transparent), image.Pt(0, 0), draw.Src)
		}

		// the canvas should be completely clean
		b := c.image.Bounds()
		test.RectIsColor(t, c.image, b.Min.X, b.Min.Y, b.Max.X, b.Max.Y, color.Transparent)
	})

	t.Run("prev regions are preserved", func(t *testing.T) {
		c := NewCanvas(320, 240)
		hs := NewHStack(c)

		square := NewRectangle(c, 20, 20, color.White)
		hs.Add(square)
		dirtyRegs, _ := hs.Render(0, 0)
		// "dirtyRegs" contains the square and only it
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 20, 20),
		}, dirtyRegs)

		hs.Remove(square)
		landscape := NewRectangle(c, 30, 10, color.White)
		hs.Add(landscape)
		dirtyRegs, _ = hs.Render(0, 0)
		// "dirtyRegs" contains the square and the landscape and only them
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 30, 10),
			image.Rect(0, 10, 20, 20),
		}, dirtyRegs)

		hs.Remove(landscape)
		portrait := NewRectangle(c, 10, 30, color.Black)
		hs.Add(portrait)
		// "dirtyRegs" contains the landscape and the portrait and only them
		// it doesn't contain the square
		dirtyRegs, _ = hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 30, 10),
			image.Rect(0, 10, 10, 30),
		}, dirtyRegs)
	})
}

// https://github.com/gchumillas/prism/issues/113
func TestIssue113(t *testing.T) {
	c := NewCanvas(320, 240)
	hs1 := NewHStack(c)
	hs2 := NewHStack(c)
	a := NewRectangle(c, 50, 50, color.Gray{Y: 75})
	b := NewRectangle(c, 50, 50, color.Gray{Y: 150})
	hs2.Add(a)
	hs2.Add(b)
	hs2.AttachTo(hs1)

	// An inner HStack contains rectangles A and B.
	// The first render returns the regions of A and B, which are cached.
	regs, _ := hs1.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(0, 0, 100, 50),
	}, regs)

	// B is removed and rendered again. The inner HStack returns only B's old
	// region, because that's the only thing that changed. Only that is cached.
	hs2.Remove(b)
	regs, _ = hs1.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(0, 0, 100, 50),
	}, regs)

	// The inner HStack is detached. prev.regs is cleared, and contains only B.
	hs2.DetachFrom(hs1)
	regs, _ = hs1.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(0, 0, 100, 50),
	}, regs)
}
