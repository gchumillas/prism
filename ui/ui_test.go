package ui_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/gchumillas/prism/internal/test"
	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"

	_ "embed"
	_ "image/png"
)

//go:embed "testdata/pen24.png"
var penPNG []byte

func TestMergeRegions(t *testing.T) {
	mirrorHor := func(rs []image.Rectangle) {
		for i, r := range rs {
			rs[i] = image.Rectangle{
				Min: image.Pt(-r.Min.X, r.Min.Y),
				Max: image.Pt(-r.Max.X, r.Max.Y),
			}.Canon()
		}
	}

	mirrorVer := func(rs []image.Rectangle) {
		for i, r := range rs {
			rs[i] = image.Rectangle{
				Min: image.Pt(r.Min.X, -r.Min.Y),
				Max: image.Pt(r.Max.X, -r.Max.Y),
			}.Canon()
		}
	}

	movements := []func([]image.Rectangle){func(r []image.Rectangle) {}, mirrorHor, mirrorVer, mirrorHor}

	t.Run("disjoint rectangles", func(t *testing.T) {
		tests := [][]image.Rectangle{
			{image.Rect(-2, -2, 2, 2), image.Rect(3, 3, 5, 5)},
			{image.Rect(-2, -2, 2, 2), image.Rect(3, 1, 5, 3)},
			{image.Rect(-2, -2, 2, 2), image.Rect(3, -1, 5, 1)},
			{image.Rect(-2, -2, 2, 2), image.Rect(3, -3, 5, 3)},
		}
		for _, test := range tests {
			for _, m := range movements {
				m(test)
				output := ui.MergeRegions(test)
				require.ElementsMatch(t, test, output)
			}
		}
	})

	t.Run("intersecting rectangles", func(t *testing.T) {
		tests := []struct {
			input    []image.Rectangle
			expected []image.Rectangle
		}{
			{
				input:    []image.Rectangle{image.Rect(-2, -2, 2, 2), image.Rect(1, 1, 3, 3)},
				expected: []image.Rectangle{image.Rect(-2, -2, 2, 1), image.Rect(-2, 1, 3, 2), image.Rect(1, 2, 3, 3)},
			},
			{
				input:    []image.Rectangle{image.Rect(-2, -2, 2, 2), image.Rect(-1, -1, 3, 1)},
				expected: []image.Rectangle{image.Rect(-2, -2, 2, -1), image.Rect(-2, -1, 3, 1), image.Rect(-2, 1, 2, 2)},
			},
		}
		for _, test := range tests {
			for _, m := range movements {
				m(test.input)
				m(test.expected)
				output := ui.MergeRegions(test.input)
				require.ElementsMatch(t, test.expected, output)
			}
		}
	})

	t.Run("nested rectangles", func(t *testing.T) {
		input := []image.Rectangle{image.Rect(-2, -2, 2, 2), image.Rect(-1, -1, 1, 1)}
		expected := []image.Rectangle{image.Rect(-2, -2, 2, 2)}

		for range 2 {
			output := ui.MergeRegions(input)
			require.Equal(t, expected, output)
			slices.Reverse(input)
			slices.Reverse(expected)
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		tests := []struct {
			input    []image.Rectangle
			expected []image.Rectangle
		}{
			{
				input:    []image.Rectangle{image.Rect(-1, -1, 1, 1), image.Rect(1, -1, 3, 1)},
				expected: []image.Rectangle{image.Rect(-1, -1, 3, 1)},
			},
			{
				input:    []image.Rectangle{image.Rect(-1, -1, 1, 1), image.Rect(1, 1, 3, 3)},
				expected: []image.Rectangle{image.Rect(-1, -1, 1, 1), image.Rect(1, 1, 3, 3)},
			},
		}
		for _, test := range tests {
			output := ui.MergeRegions(test.input)
			require.ElementsMatch(t, test.expected, output)
		}
	})

	t.Run("combined", func(t *testing.T) {
		input := []image.Rectangle{
			image.Rect(-2, -2, 2, 2),
			image.Rect(1, 0, 4, 4),
			image.Rect(2, 1, 3, 2),
			image.Rect(3, 3, 5, 5),
			image.Rect(0, 4, 1, 5),
		}
		expected := []image.Rectangle{
			image.Rect(-2, -2, 2, 0),
			image.Rect(-2, 0, 4, 2),
			image.Rect(1, 2, 4, 3),
			image.Rect(1, 3, 5, 4),
			image.Rect(0, 4, 1, 5),
			image.Rect(3, 4, 5, 5),
		}

		output := ui.MergeRegions(input)
		require.ElementsMatch(t, expected, output)
	})
}

func TestRectangle(t *testing.T) {
	t.Run("Render", func(t *testing.T) {
		app := test.NewApp()
		fillColor := color.Gray{Y: 50}
		r := app.Rectangle(50, 100, fillColor)
		r.Position.X = 100
		r.Position.Y = 50

		app.Add(r)

		app.Refresh()
		test.HlineIsColor(t, app.Screen(), 0, 100, 90, color.White)     // nothing here
		test.VlineIsColor(t, app.Screen(), 125, 0, 40, color.White)     // nothing here
		test.RectIsColor(t, app.Screen(), 0, 0, 90, 40, color.White)    // nothing here
		test.RectIsColor(t, app.Screen(), 110, 60, 145, 145, fillColor) // something here
	})
}

func TestImage(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		app := test.NewApp()
		penBuff := bytes.NewBuffer(penPNG)
		img, err := app.LoadImage(penBuff)
		img.Position.X = 30
		img.Position.Y = 30
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
		img, err := app.LoadImage(penBuff)
		img.Position.X = 10
		img.Position.Y = 20
		require.NoError(t, err)
		app.Add(img)

		// draw the image and preserve the "dirty region"
		regs := img.Render()
		require.NotEmpty(t, regs)
		r := regs[0]
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
		btn := app.LabelButton("Click me!")
		btn.Position.X = 100
		btn.Position.Y = 50
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
		btn := app.LabelButton("Click me!")
		btn.Position.X = 100
		btn.Position.Y = 50
		app.Add(btn)
		regs := btn.Render()
		require.NotEmpty(t, regs)
		r := regs[0]

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
		btn := app.LabelButton("Click me!")
		btn.Position.X = 100
		btn.Position.Y = 50
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
		btn, err := app.LoadImageButton(penBuff)
		require.NoError(t, err)
		btn.Position.X = 100
		btn.Position.Y = 50
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
		btn, err := app.LoadImageButton(penBuff)
		require.NoError(t, err)
		btn.Position.X = 100
		btn.Position.Y = 50
		app.Add(btn)
		regs := btn.Render()
		require.NotEmpty(t, regs)
		r := regs[0]

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
		btn, err := app.LoadImageButton(penBuff)
		require.NoError(t, err)
		btn.Position.X = 100
		btn.Position.Y = 50
		btn.Render()
		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})

	t.Run("Contains without rendering", func(t *testing.T) {
		app := test.NewApp()
		penBuff := bytes.NewBuffer(penPNG)
		btn, err := app.LoadImageButton(penBuff)
		require.NoError(t, err)
		require.Panics(t, func() { btn.Contains(0, 0) })
	})
}

func TestVStack(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		app := test.NewApp()
		vs := app.VStack()
		vs.Position.X = 100
		vs.Position.Y = 50
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

	t.Run("can add widgets", func(t *testing.T) {
		app := test.NewApp()

		// can stack images
		penBuff := bytes.NewBuffer(penPNG)
		imgBtn, err := app.LoadImageButton(penBuff)
		require.NoError(t, err)
		vs := app.VStack()
		vs.Add(imgBtn)

		// can stack buttons
		lbBtn := app.LabelButton("Click me!")
		vs.Add(lbBtn)

		// can stack rectangles
		r := app.Rectangle(50, 100, color.White)
		vs.Add(r)
	})

	t.Run("can remove widgets", func(t *testing.T) {
		app := test.NewApp()
		lbBtn := app.LabelButton("Click me!")
		vs := app.VStack()
		vs.Add(lbBtn)
		require.Contains(t, vs.Children(), lbBtn)
		vs.Remove(lbBtn)
		require.NotContains(t, vs.Children(), lbBtn)
	})
}
