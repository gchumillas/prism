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
	"github.com/stretchr/testify/require"
)

type Rectangle struct {
	canvas        draw.Image
	Width, Height int
	color         color.Color
}

func NewRectangle(canvas draw.Image, width, height int, c color.Color) *Rectangle {
	return &Rectangle{
		canvas: canvas,
		Width:  width,
		Height: height,
		color:  c,
	}
}

func (rec *Rectangle) Footprint() []uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, rec.Width)
	maphash.WriteComparable(&h, rec.Height)
	r, g, b, a := rec.color.RGBA()
	maphash.WriteComparable(&h, r)
	maphash.WriteComparable(&h, g)
	maphash.WriteComparable(&h, b)
	maphash.WriteComparable(&h, a)
	return []uint64{h.Sum64()}
}

func (r *Rectangle) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	rec := image.Rect(x, y, x+r.Width, y+r.Height)
	draw.Draw(r.canvas, rec, image.NewUniform(r.color), image.Pt(0, 0), draw.Src)
	return []image.Rectangle{rec}, rec.Size()
}

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

func TestHStack_Render(t *testing.T) {
	t.Run("dirty regions", func(t *testing.T) {
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
	})

	t.Run("full paint", func(t *testing.T) {
		env := NewEnv()
		recs := []*ui.Rectangle{
			ui.NewRectangle(env.ctx, 50, 150, color.Gray{Y: 100}),
			ui.NewRectangle(env.ctx, 150, 50, color.Gray{Y: 150}),
			ui.NewRectangle(env.ctx, 100, 100, color.Gray{Y: 200}),
		}

		hs := ui.NewHStack(env.ctx, ui.WithGap(10))
		for _, r := range recs {
			hs.Add(r)
		}

		regs, size := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 50, 150),
			image.Rect(60, 0, 210, 50),
			image.Rect(220, 0, 320, 100),
		}, regs)
		require.Equal(t, image.Pt(320, 150), size)
	})

	t.Run("nested items", func(t *testing.T) {
		env := NewEnv()

		r := NewRectangle(env.canvas, 50, 100, color.Gray{Y: 100})
		hs2 := ui.NewHStack(env.ctx)
		hs2.Add(r)

		hs1 := ui.NewHStack(env.ctx)
		hs1.Add(hs2)

		regs, _ := hs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 50, 100),
		}, regs)
	})

	// diff painting
	// deleting an item returns the ghost regions to the right
	t.Run("diff/remove item", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rectangle{
			NewRectangle(env.canvas, 50, 100, color.Gray{Y: 100}),
		}

		hs := ui.NewHStack(env.ctx)
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		hs.Remove(recs[0])
		regs, _ := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 50, 100),
		}, regs)
	})

	// calling Render() twice returns an empty slice
	t.Run("diff/call twice", func(t *testing.T) {
		env := NewEnv()
		recs := []*ui.Rectangle{
			ui.NewRectangle(env.ctx, 50, 150, color.Gray{Y: 100}),
			ui.NewRectangle(env.ctx, 150, 50, color.Gray{Y: 150}),
			ui.NewRectangle(env.ctx, 100, 100, color.Gray{Y: 200}),
		}

		hs := ui.NewHStack(env.ctx, ui.WithGap(10))
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		regs, _ := hs.Render(0, 0) // call Render() without changes
		require.Empty(t, regs)
	})

	// changing item height only returns the region of that item
	t.Run("diff/change item height", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rectangle{
			NewRectangle(env.canvas, 50, 150, color.Gray{Y: 100}),
			NewRectangle(env.canvas, 150, 50, color.Gray{Y: 150}),
			NewRectangle(env.canvas, 100, 100, color.Gray{Y: 200}),
		}

		hs := ui.NewHStack(env.ctx)
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		recs[1].Height += 50 // changing the height does'n affect others
		regs, _ := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(50, 0, 200, 100)}, regs)
	})

	// changing the item width returns the region of that item and the others
	t.Run("diff/change item width", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rectangle{
			NewRectangle(env.canvas, 50, 150, color.Gray{Y: 100}),
			NewRectangle(env.canvas, 150, 50, color.Gray{Y: 150}),
			NewRectangle(env.canvas, 100, 100, color.Gray{Y: 200}),
		}

		hs := ui.NewHStack(env.ctx)
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		recs[1].Width -= 50 // changing the width shift the rest of the items
		regs, _ := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(50, 0, 300, 50),
			image.Rect(150, 50, 300, 100),
		}, regs)
	})

	// deleting an item return the regions of the remaining items and the ghost region
	t.Run("diff/delete item", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rectangle{
			NewRectangle(env.canvas, 50, 150, color.Gray{Y: 100}),
			NewRectangle(env.canvas, 150, 50, color.Gray{Y: 150}),
			NewRectangle(env.canvas, 100, 100, color.Gray{Y: 200}),
		}

		hs := ui.NewHStack(env.ctx)
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		hs.Remove(recs[0]) // deleting an item shifts the rest of the items
		regs, _ := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 300, 50),
			image.Rect(0, 50, 50, 150),
			image.Rect(150, 50, 300, 100),
		}, regs)
	})

	// items with the same footprint are not repainted
	t.Run("diff/same footprint", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rectangle{
			NewRectangle(env.canvas, 50, 100, color.Gray{Y: 100}),
			NewRectangle(env.canvas, 50, 100, color.Gray{Y: 100}),
		}

		hs := ui.NewHStack(env.ctx)
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		hs.Remove(recs[0])
		regs, _ := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(50, 0, 100, 100), // only the ghost region is repainted
		}, regs)
	})

	// changing the stack position repaint all items
	t.Run("diff/change position", func(t *testing.T) {
		env := NewEnv()
		recs := []*ui.Rectangle{
			ui.NewRectangle(env.ctx, 50, 150, color.Gray{Y: 100}),
			ui.NewRectangle(env.ctx, 150, 50, color.Gray{Y: 150}),
			ui.NewRectangle(env.ctx, 100, 100, color.Gray{Y: 200}),
		}

		hs := ui.NewHStack(env.ctx)
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		regs, _ := hs.Render(100, 50)
		// cubre las regiones fantasma y las nuevas
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 300, 50),
			image.Rect(0, 50, 50, 150),
			image.Rect(100, 50, 400, 100),
			image.Rect(100, 100, 150, 200),
			image.Rect(300, 100, 400, 150),
		}, regs)
	})

	// hit area follows widget after sibling removal
	t.Run("remove button", func(t *testing.T) {
		env := NewEnv()
		hs := ui.NewHStack(env.ctx)
		btn1 := hs.AddButton(NewRectangle(env.canvas, 200, 100, color.Gray{Y: 100}))
		btn2 := hs.AddButton(NewRectangle(env.canvas, 200, 100, color.Gray{Y: 200}))

		// btn2 is clickable at position p
		p := image.Pt(300, 50)
		hs.Render(0, 0)
		require.True(t, btn2.Contains(p.X, p.Y))

		// btn2 is no longer clickable at position p
		// because removing btn1 shifts btn2 upward
		hs.RemoveButton(btn1)
		hs.Render(0, 0)
		require.False(t, btn2.Contains(p.X, p.Y))
	})
}
