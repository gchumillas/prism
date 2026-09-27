package ui_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/gchumillas/prism/internal/test"
	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"

	_ "embed"
	_ "image/png"
)

//go:embed "testdata/pen24.png"
var penPNG []byte

func TestImage(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		app := test.NewApp()
		penBuff := bytes.NewBuffer(penPNG)
		img, err := app.ImageFromReader(penBuff, ui.X(30), ui.Y(30))
		require.NoError(t, err)
		app.Add(img)

		app.Refresh()
		test.HlineIsColor(t, app.Screen(), 0, 40, 20, color.White)   // nothing here
		test.VlineIsColor(t, app.Screen(), 40, 0, 20, color.White)   // nothing here
		test.RectIsColor(t, app.Screen(), 0, 0, 20, 20, color.White) // nothing here
		test.RectHasSomeColors(t, app.Screen(), 30, 30, 54, 54)      // something here
	})

	t.Run("Render", func(t *testing.T) {
		app := test.NewApp()

		penBuff := bytes.NewBuffer(penPNG)
		img, err := app.ImageFromReader(penBuff, ui.X(10), ui.Y(20))
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
	})
}

func TestLabelButton(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		app := test.NewApp()
		btn := app.LabelButton("Click me!", ui.X(100), ui.Y(50))
		app.Add(btn)

		app.Refresh()
		test.HlineIsColor(t, app.Screen(), 0, 200, 90, color.White)  // nothing here
		test.VlineIsColor(t, app.Screen(), 200, 0, 40, color.White)  // nothing here
		test.RectIsColor(t, app.Screen(), 0, 0, 90, 40, color.White) // nothing here
		test.RectHasSomeColors(t, app.Screen(), 100, 50, 200, 100)   // something here
	})

	t.Run("Render", func(t *testing.T) {
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
	})

	t.Run("Render result", func(t *testing.T) {
		app := test.NewApp()
		btn := app.LabelButton("Click me!", ui.X(100), ui.Y(50))
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

	t.Run("Contains", func(t *testing.T) {
		app := test.NewApp()
		btn := app.LabelButton("Click me!", ui.X(100), ui.Y(50))
		btn.Render()
		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})

	t.Run("Contains without rendering", func(t *testing.T) {
		app := test.NewApp()
		btn := app.LabelButton("Click me!")
		require.Panics(t, func() { btn.Contains(0, 0) })
	})
}

func TestImageButton(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		app := test.NewApp()
		penBuff := bytes.NewBuffer(penPNG)
		btn, err := app.ImageButtonFromReader(penBuff, ui.X(100), ui.Y(50))
		require.NoError(t, err)
		app.Add(btn)

		app.Refresh()
		test.HlineIsColor(t, app.Screen(), 0, 65, 90, color.White)   // nothing here
		test.VlineIsColor(t, app.Screen(), 115, 0, 40, color.White)  // nothing here
		test.RectIsColor(t, app.Screen(), 0, 0, 90, 45, color.White) // nothing here
		test.RectHasSomeColors(t, app.Screen(), 90, 40, 140, 90)     // something here
	})

	t.Run("Render", func(t *testing.T) {
		app := test.NewApp()

		penBuff := bytes.NewBuffer(penPNG)
		btn, err := app.ImageButtonFromReader(penBuff, ui.X(100), ui.Y(50))
		require.NoError(t, err)
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

	t.Run("Contains", func(t *testing.T) {
		app := test.NewApp()

		penBuff := bytes.NewBuffer(penPNG)
		btn, err := app.ImageButtonFromReader(penBuff, ui.X(100), ui.Y(50))
		require.NoError(t, err)
		btn.Render()
		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})

	t.Run("Contains without rendering", func(t *testing.T) {
		app := test.NewApp()
		penBuff := bytes.NewBuffer(penPNG)
		btn, err := app.ImageButtonFromReader(penBuff)
		require.NoError(t, err)
		require.Panics(t, func() { btn.Contains(0, 0) })
	})
}

func TestVStack(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		app := test.NewApp()
		vs := app.VStack(ui.X(100), ui.Y(50))
		for i := range 2 {
			vs.Add(app.Label(fmt.Sprintf("Item %d", i)))
		}
		app.Add(vs)

		app.Refresh()
		test.HlineIsColor(t, app.Screen(), 0, 65, 90, color.White)   // nothing here
		test.VlineIsColor(t, app.Screen(), 115, 0, 40, color.White)  // nothing here
		test.RectIsColor(t, app.Screen(), 0, 0, 90, 45, color.White) // nothing here
		test.RectHasSomeColors(t, app.Screen(), 90, 40, 200, 125)    // something here
	})
}
