package ui_test

import (
	"image/color"
	"testing"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/internal/test"
	"github.com/gchumillas/prism/ui"
)

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
