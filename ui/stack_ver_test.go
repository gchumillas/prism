package ui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"

	_ "embed"
	_ "image/png"
)

func TestVStack_DetachFrom(t *testing.T) {
	env := NewEnv()
	hs := ui.NewHStack(env.canvas)
	vs := ui.AttachVStackTo(hs, env.canvas)
	require.True(t, vs.DetachFrom(hs))
	require.False(t, vs.DetachFrom(hs))
}

func TestVStack_AddButton(t *testing.T) {
	t.Run("add button", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		btn := ui.AttachButtonTo(vs, NewRect(env.image, 0, 0, 100, 50, color.White))
		vs.Render(0, 0)
		require.True(t, btn.Contains(50, 25))
	})

	t.Run("repaint stacker in a different position", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		btn := ui.AttachButtonTo(vs, NewRect(env.image, 0, 0, 100, 50, color.White))
		vs.Render(0, 0)
		vs.Render(100, 50)
		require.True(t, btn.Contains(150, 75))
	})

	t.Run("remove item", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		r := NewRect(env.image, 0, 0, 100, 50, color.Black)
		vs.Add(r)
		btn := ui.AttachButtonTo(vs, NewRect(env.image, 0, 0, 100, 50, color.White))
		vs.Render(0, 0)
		require.True(t, btn.Contains(50, 75))

		// When an item is removed, the button moves up.
		// Verify that it is clickable in its new position and not in its old position.
		vs.Remove(r)
		vs.Render(0, 0)
		require.True(t, btn.Contains(50, 25))
		require.False(t, btn.Contains(50, 75))
	})

	t.Run("change item height", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		r := NewRect(env.image, 0, 0, 100, 50, color.Black)
		vs.Add(r)
		btn := ui.AttachButtonTo(vs, NewRect(env.image, 0, 0, 100, 50, color.White))
		vs.Render(0, 0)
		require.True(t, btn.Contains(50, 75))

		// When an item's height changes, the button moves vertically.
		// Verify that it is clickable in its new position and not in its old position.
		r.Max.Y += 50
		vs.Render(0, 0)
		require.True(t, btn.Contains(50, 125))
		require.False(t, btn.Contains(50, 75))
	})
}

func TestVStack_Render(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas, ui.WithGap(10))
		vs.Add(NewRect(env.image, 0, 0, 100, 50, color.Black))
		vs.Add(NewRect(env.image, 0, 0, 50, 50, color.White))
		_, size := vs.Render(0, 0)
		require.Equal(t, image.Pt(100, 110), size)
	})

	t.Run("render all widgets", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.canvas)
		for i := range 3 {
			shape := NewRect(env.image, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			vs.Add(shape)
		}

		// all shapes are rendered
		regs, _ := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 150)}, regs)

		// no shape is rendered
		regs, _ = vs.Render(0, 0)
		require.Empty(t, regs)
	})

	t.Run("render only affected widget", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.canvas)
		var shapes []*Rect
		for i := range 3 {
			shape := NewRect(env.image, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		shapes[1].Max.X += 50
		regs, _ := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 50, 150, 100)}, regs)
	})

	t.Run("render affected widget and below", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.canvas)
		var shapes []*Rect
		for i := range 3 {
			shape := NewRect(env.image, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		shapes[1].Max.Y += 50
		regs, _ := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 50, 100, 200)}, regs)
	})

	t.Run("render widgets below deleted", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.canvas)
		var shapes []*Rect
		for i := range 3 {
			shape := NewRect(env.image, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		vs.Remove(shapes[0])
		regs, _ := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 150)}, regs)
	})

	t.Run("render on change position", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.canvas)
		for i := range 3 {
			shape := NewRect(env.image, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			vs.Add(shape)
		}

		vs.Render(0, 0)
		regs, _ := vs.Render(100, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 200, 150)}, regs)
	})

	t.Run("nested", func(t *testing.T) {
		env := NewEnv()

		vs1 := ui.NewVStack(env.canvas)
		vs2 := ui.AttachVStackTo(vs1, env.canvas)

		sh := NewRect(env.image, 0, 0, 100, 50, color.White)
		vs2.Add(sh)

		regs, _ := vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 50)}, regs)

		sh.Max.X += 50
		regs, _ = vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 150, 50)}, regs)
	})

	t.Run("nested change after steady state", func(t *testing.T) {
		env := NewEnv()

		vs1 := ui.NewVStack(env.canvas)
		vs2 := ui.AttachVStackTo(vs1, env.canvas)

		sh := NewRect(env.image, 0, 0, 100, 50, color.White)
		vs2.Add(sh)

		regs, _ := vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 50)}, regs)

		// render until no widget changes (steady state)
		vs1.Render(0, 0)
		vs1.Render(0, 0)

		// a nested change must still be detected by the outer stack
		sh.Max.X += 50
		regs, _ = vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 150, 50)}, regs)
	})

	t.Run("hit area follows widget after sibling removal", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		btn1 := ui.AttachLabelButtonTo(vs, env.canvas, "Button")
		btn2 := ui.AttachLabelButtonTo(vs, env.canvas, "Button")

		// btn2 is clickable at position p
		p := image.Pt(80, 115)
		vs.Render(0, 0)
		require.True(t, btn2.Contains(p.X, p.Y))

		// btn2 is no longer clickable at position p
		// because removing btn1 shifts btn2 upward
		btn1.DetachFrom(vs)
		vs.Render(0, 0)
		require.False(t, btn2.Contains(p.X, p.Y))
	})

	t.Run("remove from list of different widgets", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		var shapes []*Rect
		for i := range uint8(2) {
			r := NewRect(env.image, 0, 0, 100, 50, color.Gray{Y: i})
			shapes = append(shapes, r)
			vs.Add(r)
		}
		vs.Render(0, 0)

		sh0 := shapes[0]
		vs.Remove(sh0)
		regs, _ := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 100)}, regs)
	})

	t.Run("remove from list of identical widgets", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		var shapes []*Rect
		for range uint8(3) {
			r := NewRect(env.image, 0, 0, 100, 50, color.White)
			shapes = append(shapes, r)
			vs.Add(r)
		}
		vs.Render(0, 0)

		sh0 := shapes[0]
		vs.Remove(sh0)
		regs, _ := vs.Render(0, 0)
		// only returns the region of the last rectangle
		// since the first rectangle is identical to the second and the third,
		// there is no need to repaint those regions
		require.Equal(t, []image.Rectangle{image.Rect(0, 100, 100, 150)}, regs)
	})

	t.Run("render different widget sizes", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rect{
			NewRect(env.image, 0, 0, 100, 50, color.Gray{Y: 25}),  // landscape
			NewRect(env.image, 0, 50, 50, 100, color.Gray{Y: 50}), // square
		}
		vs := ui.NewVStack(env.canvas)
		for _, rec := range recs {
			vs.Add(rec)
		}
		vs.Render(0, 0)
		vs.Remove(recs[0])
		regs, _ := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 50),
			image.Rect(0, 50, 50, 100),
		}, regs)
	})

	t.Run("different widget with the same footprint", func(t *testing.T) {
		env := NewEnv()
		recs := []ui.Widget{
			NewRect(env.image, 0, 0, 100, 50, color.White),
			NewRect1(env.image, 0, 50, 100, 100, color.White),
		}
		vs := ui.NewVStack(env.canvas, ui.WithGap(10))
		for _, rec := range recs {
			vs.Add(rec)
		}
		regs, _ := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 50),
			image.Rect(0, 60, 100, 110),
		}, regs)

		vs.Remove(recs[0])
		regs, _ = vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 50),   // must be repainted, since Rect1 and Rect are of different types
			image.Rect(0, 60, 100, 110), // and the abandoned region must also be repainted
		}, regs)
	})

	t.Run("nested vstack child replaced by different type", func(t *testing.T) {
		env := NewEnv()

		vs1 := ui.NewVStack(env.canvas)
		vs2 := ui.AttachVStackTo(vs1, env.canvas)

		sh := NewRect(env.image, 0, 0, 100, 50, color.Black)
		vs2.Add(sh)
		regs, _ := vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 50)}, regs)

		// replace child with another one with identical footprint and different type
		vs2.Remove(sh)
		vs2.Add(NewRect1(env.image, 0, 0, 100, 50, color.Black))
		regs, _ = vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 50)}, regs)
	})
}
