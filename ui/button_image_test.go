package ui_test

import (
	"bytes"
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/gchumillas/prism/internal/test"
	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestImageButton(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight, WithPadding(image.Pt(100, 50)))

		dat, err := os.ReadFile("testdata/icon3.png")
		require.NoError(t, err)
		btn, err := ui.LoadImageButton(ctx, bytes.NewReader(dat))
		require.NoError(t, err)
		app.Add(btn)
		app.Render()

		s.Refresh()
		test.RectBorderIsColor(t, s.Display, 95, 45, 140, 90, color.White)
		test.DotHasSomeColors(t, s.Display, 100, 50, 10)       // top-left corner
		test.DotHasSomeColors(t, s.Display, 135, 50, 10)       // top-right corner
		test.DotHasSomeColors(t, s.Display, 100, 85, 10)       // bottom-left corner
		test.DotHasSomeColors(t, s.Display, 135, 85, 10)       // bottom-right corner
		test.RectHasSomeColors(t, s.Display, 105, 55, 125, 75) // icon
	})

	t.Run("Render", func(t *testing.T) {
		padding := image.Pt(100, 50)
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight, WithPadding(padding))

		dat, err := os.ReadFile("testdata/icon1.png")
		require.NoError(t, err)
		btn, err := ui.LoadImageButton(ctx, bytes.NewReader(dat))
		require.NoError(t, err)
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

		dat, err := os.ReadFile("testdata/icon1.png")
		require.NoError(t, err)
		btn, err := ui.LoadImageButton(ctx, bytes.NewReader(dat))
		require.NoError(t, err)
		btn.Render(100, 50)

		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})
}
