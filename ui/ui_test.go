package ui_test

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/internal/test"
	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"

	_ "embed"
	_ "image/png"
)

//go:embed "testdata/pen24.png"
var penPNG []byte

func TestScene_ImageFromReader(t *testing.T) {
	canvas := gg.NewContext(30, 30)
	s := ui.NewScene(canvas)

	penBuff := bytes.NewBuffer(penPNG)
	img, err := s.ImageFromReader(penBuff)
	require.NoError(t, err)
	img.Render()

	snap := test.GenSnapshot(t, canvas)
	test.RectHasSomeColors(t, snap, 0, 0, 24, 24)
}

func TestImage_Render(t *testing.T) {
	canvas := gg.NewContext(30, 30)
	s := ui.NewScene(canvas)

	penBuff := bytes.NewBuffer(penPNG)
	img, err := s.ImageFromReader(penBuff)
	require.NoError(t, err)
	r := img.Render()
	require.Equal(t, image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 24, Y: 24},
	}, r)
}

func TestLabelButton_Render(t *testing.T) {
	canvas := gg.NewContext(220, 80)
	s := ui.NewScene(canvas)
	btn := s.LabelButton("Click me!")
	btn.Render()

	snap := test.GenSnapshot(t, canvas)
	test.VlineIsColor(t, snap, 100, 5, 10, color.Transparent) // top padding
	test.HlineIsColor(t, snap, 10, 35, 10, color.Transparent) // left padding
	test.RectHasSomeColors(t, snap, 35, 20, 170, 50)          // text
	test.DotHasSomeColors(t, snap, 0, 0, 5)                   // left-top corner
	test.DotHasSomeColors(t, snap, 205, 0, 5)                 // right-top corner
	test.DotHasSomeColors(t, snap, 205, 75, 5)                // right-bottom corner
	test.DotHasSomeColors(t, snap, 0, 75, 5)                  // left-bottom corner
}
