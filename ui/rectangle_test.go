package ui_test

import (
	"image/color"
	"testing"

	"github.com/gchumillas/prism/internal/test"
	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestRectangle_Render(t *testing.T) {
	t.Run("change fill color", func(t *testing.T) {
		env := NewEnv()
		rec := ui.NewRectangle(env.ctx, 100, 50, color.White)
		rec.Render(0, 0)

		canvas := env.ctx.Canvas()
		c := canvas.At(10, 10)
		require.Equal(t, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, c)

		rec.FillColor = color.Black
		rec.Render(0, 0)
		c = canvas.At(10, 10)
		require.Equal(t, color.RGBA{A: 0xff}, c)
	})

	t.Run("change size", func(t *testing.T) {
		env := NewEnv()
		rec := ui.NewRectangle(env.ctx, 100, 50, color.White)
		rec.Render(10, 10)

		test.RectBorderIsColor(t, env.canvas, 0, 0, 120, 70, color.RGBA{}) // border around
		test.RectHasSomeColors(t, env.canvas, 0, 0, 120, 70)               // rect covering
		test.RectIsColor(t, env.canvas, 15, 15, 105, 55, color.White)      // rect inside

		rec.Width, rec.Height = 200, 100
		rec.Render(10, 10)
		test.RectBorderIsColor(t, env.canvas, 0, 0, 220, 120, color.RGBA{}) // border around
		test.RectHasSomeColors(t, env.canvas, 0, 0, 220, 120)               // rect covering
		test.RectIsColor(t, env.canvas, 15, 15, 205, 105, color.White)      // rect inside
	})

	t.Run("can be stacked vertically", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)

		r := ui.NewRectangle(env.ctx, 50, 100, color.White)
		vs.Add(r)
		require.Contains(t, vs.Children(), r)
	})
}
