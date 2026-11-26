package ui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestHStack_DetachFrom(t *testing.T) {
	vs := ui.NewVStack()
	hs := ui.NewHStack()
	_ = hs.AttachTo(vs)
	require.NoError(t, hs.DetachFrom(vs))
	require.ErrorIs(t, hs.DetachFrom(vs), ui.ErrNotAttached)
}

func TestHStack_Render(t *testing.T) {
	t.Run("full paint", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rectangle{
			NewRectangle(env.canvas, 50, 150, color.Gray{Y: 100}),
			NewRectangle(env.canvas, 150, 50, color.Gray{Y: 150}),
			NewRectangle(env.canvas, 100, 100, color.Gray{Y: 200}),
		}

		hs := ui.NewHStack(ui.WithGap(10))
		for _, r := range recs {
			hs.Add(r)
		}

		regs := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 50, 150),
			image.Rect(60, 0, 210, 50),
			image.Rect(220, 0, 320, 100),
		}, regs)
		size := hs.Size()
		require.Equal(t, image.Pt(320, 150), size)
	})

	t.Run("nested items", func(t *testing.T) {
		env := NewEnv()

		r := NewRectangle(env.canvas, 50, 100, color.Gray{Y: 100})
		hs1 := ui.NewHStack()
		hs2 := ui.NewHStack()
		_ = hs2.AttachTo(hs1)
		hs2.Add(r)

		regs := hs1.Render(0, 0)
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

		hs := ui.NewHStack()
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		hs.Remove(recs[0])
		regs := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 50, 100),
		}, regs)
	})

	// calling Render() twice returns an empty slice
	t.Run("diff/call twice", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rectangle{
			NewRectangle(env.canvas, 50, 150, color.Gray{Y: 100}),
			NewRectangle(env.canvas, 150, 50, color.Gray{Y: 150}),
			NewRectangle(env.canvas, 100, 100, color.Gray{Y: 200}),
		}

		hs := ui.NewHStack(ui.WithGap(10))
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		regs := hs.Render(0, 0) // call Render() without changes
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

		hs := ui.NewHStack()
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		recs[1].height += 50 // changing the height does'n affect others
		regs := hs.Render(0, 0)
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

		hs := ui.NewHStack()
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		recs[1].width -= 50 // changing the width shift the rest of the items
		regs := hs.Render(0, 0)
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

		hs := ui.NewHStack()
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		hs.Remove(recs[0]) // deleting an item shifts the rest of the items
		regs := hs.Render(0, 0)
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

		hs := ui.NewHStack()
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		hs.Remove(recs[0])
		regs := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(50, 0, 100, 100), // only the ghost region is repainted
		}, regs)
	})

	// changing the stack position repaint all items
	t.Run("diff/change position", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rectangle{
			NewRectangle(env.canvas, 50, 150, color.Gray{Y: 100}),
			NewRectangle(env.canvas, 150, 50, color.Gray{Y: 150}),
			NewRectangle(env.canvas, 100, 100, color.Gray{Y: 200}),
		}

		hs := ui.NewHStack()
		for _, r := range recs {
			hs.Add(r)
		}

		hs.Render(0, 0)
		regs := hs.Render(100, 50)
		// covers the ghost regions and the new ones
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
		hs := ui.NewHStack()
		btn1 := ui.NewButton(NewRectangle(env.canvas, 200, 100, color.Gray{Y: 100}))
		btn2 := ui.NewButton(NewRectangle(env.canvas, 200, 100, color.Gray{Y: 200}))
		_ = btn1.AttachTo(hs)
		_ = btn2.AttachTo(hs)

		// btn2 is clickable at position p
		p := image.Pt(300, 50)
		hs.Render(0, 0)
		require.True(t, btn2.Contains(p.X, p.Y))

		// btn2 is no longer clickable at position p
		// because removing btn1 shifts btn2 upward
		_ = btn1.DetachFrom(hs)
		hs.Render(0, 0)
		require.False(t, btn2.Contains(p.X, p.Y))
	})
}
