package ui_test

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/gchumillas/prism/internal/test"
	"github.com/stretchr/testify/require"

	_ "embed"
	_ "image/png"
)

//go:embed "testdata/pen24.png"
var penPNG []byte

func TestContext_ImageFromReader(t *testing.T) {
	app := test.NewApp()

	penBuff := bytes.NewBuffer(penPNG)
	img, err := app.ImageFromReader(penBuff)
	require.NoError(t, err)
	app.Add(img)

	app.Refresh()
	test.RectHasSomeColors(t, app.Screen(), 0, 0, 24, 24)
}

func TestImage_Render(t *testing.T) {
	app := test.NewApp()

	penBuff := bytes.NewBuffer(penPNG)
	img, err := app.ImageFromReader(penBuff)
	img.SetX(10)
	img.SetY(20)
	require.NoError(t, err)
	app.Add(img)

	// draw the image and preserve the "dirty region"
	r := img.Render()
	require.Equal(t, image.Rect(10, 20, 34, 44), r)
	app.Refresh()

	// check the element is there
	test.RectHasSomeColors(t, app.Screen(), r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)

	// remove the element and refresh only "dirty region"
	app.Remove(img)
	app.RefreshRect(r)

	// check the element is NOT there
	test.RectIsColor(t, app.Screen(), r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color.White)
}

func TestLabelButton_Render(t *testing.T) {
	app := test.NewApp()
	btn := app.LabelButton("Click me!")
	app.Add(btn)

	app.Refresh()
	screen := app.Screen()
	test.VlineIsColor(t, screen, 100, 5, 10, color.White) // top padding
	test.HlineIsColor(t, screen, 10, 35, 10, color.White) // left padding
	test.RectHasSomeColors(t, screen, 35, 20, 170, 50)    // text
	test.DotHasSomeColors(t, screen, 0, 0, 5)             // left-top corner
	test.DotHasSomeColors(t, screen, 205, 0, 5)           // right-top corner
	test.DotHasSomeColors(t, screen, 205, 75, 5)          // right-bottom corner
	test.DotHasSomeColors(t, screen, 0, 75, 5)            // left-bottom corner

	t.Run("dirty rect", func(t *testing.T) {
		app := test.NewApp()
		btn := app.LabelButton("Click me!")
		btn.SetX(100)
		btn.SetY(50)
		app.Add(btn)
		r := btn.Render()

		// the widget is there
		app.Refresh()
		test.RectHasSomeColors(t, app.Screen(), r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)

		// the widget is NOT there
		app.Remove(btn)
		app.RefreshRect(r)
		test.RectIsColor(t, app.Screen(), r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color.White)
	})
}

func TestLabelButton_Contains(t *testing.T) {
	app := test.NewApp()
	btn := app.LabelButton("Click me!")
	btn.SetX(100)
	btn.SetY(50)
	btn.Render()
	require.False(t, btn.Contains(0, 0))   // outside
	require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
	require.True(t, btn.Contains(110, 60)) // inside

	t.Run("forgot calling Render", func(t *testing.T) {
		app := test.NewApp()
		btn := app.LabelButton("Click me!")
		require.Panics(t, func() { btn.Contains(0, 0) })
	})
}

func TestImageButton_Render(t *testing.T) {
	app := test.NewApp()

	penBuff := bytes.NewBuffer(penPNG)
	btn, err := app.ImageButtonFromReader(penBuff)
	require.NoError(t, err)
	btn.SetX(100)
	btn.SetY(50)
	app.Add(btn)
	r := btn.Render()

	// the widget is there
	app.Refresh()
	test.RectHasSomeColors(t, app.Screen(), r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)

	// the widget is NOT there
	app.Remove(btn)
	app.RefreshRect(r)
	test.RectIsColor(t, app.Screen(), r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color.White)
}

func TestImageButton_Contains(t *testing.T) {
	app := test.NewApp()

	penBuff := bytes.NewBuffer(penPNG)
	btn, err := app.ImageButtonFromReader(penBuff)
	require.NoError(t, err)
	btn.SetX(100)
	btn.SetY(50)
	btn.Render()
	require.False(t, btn.Contains(0, 0))   // outside
	require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
	require.True(t, btn.Contains(110, 60)) // inside

	t.Run("forgot calling Render", func(t *testing.T) {
		app := test.NewApp()
		penBuff := bytes.NewBuffer(penPNG)
		btn, err := app.ImageButtonFromReader(penBuff)
		require.NoError(t, err)
		require.Panics(t, func() { btn.Contains(0, 0) })
	})
}
