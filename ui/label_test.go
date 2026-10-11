package ui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/gchumillas/prism/internal/test"
	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestLabelButton(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		padding := image.Pt(100, 50)
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight, WithPadding(padding))

		btn := ui.NewLabelButton(ctx, "Click me!")
		app.Add(btn)
		app.Render()

		s.Refresh()
		test.RectBorderIsColor(t, s.Display, 90, 40, 315, 135, color.White)
		test.DotHasSomeColors(t, s.Display, 100, 50, 5)         // top-left corner
		test.DotHasSomeColors(t, s.Display, 307, 50, 5)         // top-right corner
		test.DotHasSomeColors(t, s.Display, 100, 125, 5)        // bottom-left corner
		test.DotHasSomeColors(t, s.Display, 305, 125, 5)        // bottom-right corner
		test.RectHasSomeColors(t, s.Display, 130, 70, 275, 100) // label
	})

	t.Run("Render", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)

		btn := ui.NewLabelButton(ctx, "Click me!")
		app.Add(btn)
		app.Render()

		s.Refresh()
		test.VlineIsColor(t, s.Display, 100, 5, 10, color.White) // top padding
		test.HlineIsColor(t, s.Display, 10, 35, 10, color.White) // left padding
		test.RectHasSomeColors(t, s.Display, 35, 20, 170, 50)    // text
		test.DotHasSomeColors(t, s.Display, 0, 0, 5)             // left-top corner
		test.DotHasSomeColors(t, s.Display, 205, 0, 5)           // right-top corner
		test.DotHasSomeColors(t, s.Display, 205, 75, 5)          // right-bottom corner
		test.DotHasSomeColors(t, s.Display, 0, 75, 5)            // left-bottom corner
	})

	t.Run("Render result", func(t *testing.T) {
		padding := image.Pt(100, 50)
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight, WithPadding(padding))

		btn := ui.NewLabelButton(ctx, "Click me!")
		app.Add(btn)
		regs := btn.Render(padding.X, padding.Y)
		require.NotEmpty(t, regs)
		r := regs[0]

		// the widget is there
		app.Render()
		s.Refresh()
		test.RectHasSomeColors(t, s.Display, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)

		// the widget is NOT there
		app.Remove(btn)
		app.Render()
		s.RefreshRecs([]image.Rectangle{r})
		test.RectIsColor(t, s.Display, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color.White)
	})

	t.Run("Contains", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)

		btn := ui.NewLabelButton(ctx, "Click me!")
		btn.Render(100, 50)

		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})
}
