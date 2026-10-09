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

func TestImage(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(30, 30)

		dat, err := os.ReadFile("testdata/icon1.png")
		require.NoError(t, err)
		img, err := ui.LoadImage(ctx, bytes.NewReader(dat))
		require.NoError(t, err)
		app.Add(img)
		app.Render()

		s.Refresh()
		test.RectBorderIsColor(t, s.Display, 25, 25, 55, 55, color.White)
		test.RectHasSomeColors(t, s.Display, 30, 30, 50, 50)
	})

	t.Run("Render", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(10, 20)

		dat, err := os.ReadFile("testdata/icon2.png")
		require.NoError(t, err)
		img, err := ui.LoadImage(ctx, bytes.NewReader(dat))
		require.NoError(t, err)
		app.Add(img)

		// draw the image and preserve the "dirty region"
		regs := img.Render(app.Padding.X, app.Padding.Y)
		require.NotEmpty(t, regs)
		r := regs[0]
		require.Equal(t, image.Rect(10, 20, 34, 44), r)

		// check the element is there
		app.Render()
		s.Refresh()
		test.RectHasSomeColors(t, s.Display, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)

		// remove the element and refresh only "dirty region"
		app.Remove(img)
		app.Render()
		s.RefreshRecs([]image.Rectangle{r})

		// check the element is NOT there
		test.RectIsColor(t, s.Display, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color.White)
	})
}
