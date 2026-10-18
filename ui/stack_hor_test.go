package ui_test

import (
	"hash/maphash"
	"image"
	"image/color"
	"image/draw"
	"testing"

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

func TestHStack_DetachFrom(t *testing.T) {
	env := NewEnv()
	vs := ui.NewVStack(env.ctx)
	hs := ui.AttachHStackTo(vs, env.ctx)
	require.True(t, hs.DetachFrom(vs))
	require.False(t, hs.DetachFrom(vs))
}

func TestHStack_Render(t *testing.T) {
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
		hs1 := ui.NewHStack(env.ctx)
		hs2 := ui.AttachHStackTo(hs1, env.ctx)
		hs2.Add(r)

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
		btn1 := ui.AttachButtonTo(hs, NewRectangle(env.canvas, 200, 100, color.Gray{Y: 100}))
		btn2 := ui.AttachButtonTo(hs, NewRectangle(env.canvas, 200, 100, color.Gray{Y: 200}))

		// btn2 is clickable at position p
		p := image.Pt(300, 50)
		hs.Render(0, 0)
		require.True(t, btn2.Contains(p.X, p.Y))

		// btn2 is no longer clickable at position p
		// because removing btn1 shifts btn2 upward
		btn1.DetachFrom(hs)
		hs.Render(0, 0)
		require.False(t, btn2.Contains(p.X, p.Y))
	})
}
