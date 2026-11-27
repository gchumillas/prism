package ui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestText_Footprint(t *testing.T) {
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))

	t.Run("color", func(t *testing.T) {
		t1 := ui.NewText(c, "lorem ipsum", 200)
		t2 := ui.NewText(c, "lorem ipsum", 200)
		t1.Color = color.Black
		t2.Color = color.Transparent
		require.NotEqual(t, t1.Footprint(), t2.Footprint())
	})

	t.Run("font size", func(t *testing.T) {
		t1 := ui.NewText(c, "lorem ipsum", 200, ui.WithFontSize(12))
		t2 := ui.NewText(c, "lorem ipsum", 200, ui.WithFontSize(12.5))
		require.NotEqual(t, t1.Footprint(), t2.Footprint())
	})

	t.Run("texts", func(t *testing.T) {
		l1 := ui.NewText(c, "lorem ipsum", 200)
		l2 := ui.NewText(c, "lorem ipsum dolor", 200)
		require.NotEqual(t, l1.Footprint(), l2.Footprint())
	})

	t.Run("texts", func(t *testing.T) {
		l1 := ui.NewText(c, "lorem ipsum", 200)
		l2 := ui.NewText(c, "lorem ipsum", 300)
		require.NotEqual(t, l1.Footprint(), l2.Footprint())
	})
}
