package ui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestLabel_Render(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		env := NewEnv()
		lb := ui.NewLabel(env.canvas, "Hello, World!", 10)
		size := lb.Size()
		require.Equal(t, image.Pt(181, 37), size)
	})
}

func TestLabel_Footprint(t *testing.T) {
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))

	t.Run("color", func(t *testing.T) {
		l1 := ui.NewLabel(c, "label", 12)
		l2 := ui.NewLabel(c, "label", 12)
		l1.Color = color.Black
		l2.Color = color.Transparent
		require.NotEqual(t, l1.Footprint(), l2.Footprint())
	})

	t.Run("font size", func(t *testing.T) {
		l1 := ui.NewLabel(c, "label", 12)
		l2 := ui.NewLabel(c, "label", 12.5)
		require.NotEqual(t, l1.Footprint(), l2.Footprint())
	})

	t.Run("texts", func(t *testing.T) {
		l1 := ui.NewLabel(c, "label 1", 12)
		l2 := ui.NewLabel(c, "label 2", 12)
		require.NotEqual(t, l1.Footprint(), l2.Footprint())
	})
}
