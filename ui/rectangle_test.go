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
		ctx := ui.NewContext(320, 240)
		rec := ui.NewRectangle(ctx, 100, 50, color.White)
		rec.Render(0, 0)

		canvas := ctx.Canvas()
		c := canvas.At(10, 10)
		require.Equal(t, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, c)

		rec.FillColor = color.Black
		rec.Render(0, 0)
		c = canvas.At(10, 10)
		require.Equal(t, color.RGBA{A: 0xff}, c)
	})

	t.Run("change size", func(t *testing.T) {
		ctx := ui.NewContext(320, 240)
		rec := ui.NewRectangle(ctx, 100, 50, color.White)
		rec.Render(10, 10)

		canvas := ctx.Canvas()
		test.RectBorderIsColor(t, canvas, 0, 0, 120, 70, color.RGBA{}) // border around
		test.RectHasSomeColors(t, canvas, 0, 0, 120, 70)               // rect covering
		test.RectIsColor(t, canvas, 15, 15, 105, 55, color.White)      // rect inside

		rec.Width, rec.Height = 200, 100
		rec.Render(10, 10)
		test.RectBorderIsColor(t, canvas, 0, 0, 220, 120, color.RGBA{}) // border around
		test.RectHasSomeColors(t, canvas, 0, 0, 220, 120)               // rect covering
		test.RectIsColor(t, canvas, 15, 15, 205, 105, color.White)      // rect inside
	})
}
