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

const (
	appWidth  = 640
	appHeight = 480
)

type App struct {
	*ui.VStack
	bg *ui.Rectangle
}

func NewApp(ctx ui.Context, w, h int) *App {
	return &App{
		VStack: ui.NewVStack(ctx),
		bg:     ui.NewRectangle(ctx, w, h, color.White),
	}
}

func (s *App) Render() {
	s.bg.Render(0, 0)
	s.VStack.Render(0, 0)
}

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
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(100, 50)

		fillColor := color.Gray{Y: 50}
		r := ui.NewRectangle(ctx, 50, 100, fillColor)
		app.Add(r)
		app.Render()

		s.Refresh()
		test.HlineIsColor(t, s.Display, 0, 100, 90, color.White)     // nothing here
		test.VlineIsColor(t, s.Display, 125, 0, 40, color.White)     // nothing here
		test.RectIsColor(t, s.Display, 0, 0, 90, 40, color.White)    // nothing here
		test.RectIsColor(t, s.Display, 110, 60, 145, 145, fillColor) // something here
	})
}

func TestImage(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(30, 30)

		penBuff := bytes.NewBuffer(penPNG)
		img, err := ui.LoadImage(ctx, penBuff)
		require.NoError(t, err)
		app.Add(img)
		app.Render()

		s.Refresh()
		test.HlineIsColor(t, s.Display, 0, 40, 20, color.White)   // nothing here
		test.VlineIsColor(t, s.Display, 40, 0, 20, color.White)   // nothing here
		test.RectIsColor(t, s.Display, 0, 0, 20, 20, color.White) // nothing here
		test.RectHasSomeColors(t, s.Display, 30, 30, 54, 54)      // something here
	})

	t.Run("Render", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(10, 20)

		penBuff := bytes.NewBuffer(penPNG)
		img, err := ui.LoadImage(ctx, penBuff)
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

func TestLabelButton(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(100, 50)

		btn := ui.NewLabelButton(ctx, "Click me!")
		app.Add(btn)
		app.Render()

		s.Refresh()
		test.HlineIsColor(t, s.Display, 0, 200, 90, color.White)  // nothing here
		test.VlineIsColor(t, s.Display, 200, 0, 40, color.White)  // nothing here
		test.RectIsColor(t, s.Display, 0, 0, 90, 40, color.White) // nothing here
		test.RectHasSomeColors(t, s.Display, 100, 50, 200, 100)   // something here
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
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(100, 50)

		btn := ui.NewLabelButton(ctx, "Click me!")
		app.Add(btn)
		regs := btn.Render(app.Padding.X, app.Padding.Y)
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

	t.Run("Contains without rendering", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)
		btn := ui.NewLabelButton(ctx, "Click me!")
		require.PanicsWithError(t, ui.ErrEmptyHitbox.Error(), func() {
			btn.Contains(0, 0)
		})
	})
}

func TestImageButton(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(100, 50)

		penBuff := bytes.NewBuffer(penPNG)
		btn, err := ui.LoadImageButton(ctx, penBuff)
		require.NoError(t, err)
		app.Add(btn)
		app.Render()

		s.Refresh()
		test.HlineIsColor(t, s.Display, 0, 65, 90, color.White)   // nothing here
		test.VlineIsColor(t, s.Display, 115, 0, 40, color.White)  // nothing here
		test.RectIsColor(t, s.Display, 0, 0, 90, 45, color.White) // nothing here
		test.RectHasSomeColors(t, s.Display, 90, 40, 140, 90)     // something here
	})

	t.Run("Render", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(100, 50)

		penBuff := bytes.NewBuffer(penPNG)
		btn, err := ui.LoadImageButton(ctx, penBuff)
		require.NoError(t, err)
		app.Add(btn)
		regs := btn.Render(app.Padding.X, app.Padding.Y)
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

		penBuff := bytes.NewBuffer(penPNG)
		btn, err := ui.LoadImageButton(ctx, penBuff)
		require.NoError(t, err)
		btn.Render(100, 50)

		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})

	t.Run("Contains without rendering", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)
		penBuff := bytes.NewBuffer(penPNG)
		btn, err := ui.LoadImageButton(ctx, penBuff)
		require.NoError(t, err)
		require.PanicsWithError(t, ui.ErrEmptyHitbox.Error(), func() {
			btn.Contains(0, 0)
		})
	})
}

type CustomWidget struct{ Label ui.Label }

func NewCustomWidget(c ui.Context, text string) *CustomWidget {
	return &CustomWidget{Label: *ui.NewLabel(c, text)}
}
func (cw *CustomWidget) Render(x, y int) []image.Rectangle { return cw.Label.Render(x, y) }
func (cw *CustomWidget) Size() image.Point                 { return cw.Label.Size() }

func TestVStack(t *testing.T) {
	t.Run("X Y options", func(t *testing.T) {
		s := test.NewScreen(appWidth, appHeight)
		ctx := ui.NewContextFromRGBA(s.Canvas)
		app := NewApp(ctx, appWidth, appHeight)
		app.Padding = image.Pt(100, 50)

		vs := ui.NewVStack(ctx)
		for i := range 2 {
			vs.Add(ui.NewLabel(ctx, fmt.Sprintf("Item %d", i)))
		}
		app.Add(vs)
		app.Render()

		s.Refresh()
		test.HlineIsColor(t, s.Display, 0, 65, 90, color.White)   // nothing here
		test.VlineIsColor(t, s.Display, 115, 0, 40, color.White)  // nothing here
		test.RectIsColor(t, s.Display, 0, 0, 90, 45, color.White) // nothing here
		test.RectHasSomeColors(t, s.Display, 90, 40, 200, 125)    // something here
	})

	t.Run("can add widgets", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)
		vs := ui.NewVStack(ctx)

		// can stack images
		penBuff := bytes.NewBuffer(penPNG)
		imgBtn, err := ui.LoadImageButton(ctx, penBuff)
		require.NoError(t, err)
		vs.Add(imgBtn)

		// can stack buttons
		lbBtn := ui.NewLabelButton(ctx, "Click me!")
		vs.Add(lbBtn)

		// can stack rectangles
		r := ui.NewRectangle(ctx, 50, 100, color.White)
		vs.Add(r)
	})

	t.Run("can remove widgets", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)
		vs := ui.NewVStack(ctx)

		lbBtn := ui.NewLabelButton(ctx, "Click me!")
		vs.Add(lbBtn)
		require.Contains(t, vs.Children(), lbBtn)
		vs.Remove(lbBtn)
		require.NotContains(t, vs.Children(), lbBtn)
	})

	t.Run("can create stackable widgets", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)
		vs := ui.NewVStack(ctx)
		cw := NewCustomWidget(ctx, "custom ui")
		vs.Add(cw)
	})
}
