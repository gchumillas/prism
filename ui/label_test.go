package ui_test

import (
	"image"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestLabel_Render(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		env := NewEnv()
		lb := ui.NewLabel(env.canvas, "Hello, World!")
		_, size := lb.Render(0, 0)
		require.Equal(t, image.Pt(181, 37), size)
	})
}

func TestLabel_Footprint(t *testing.T) {
	c := ui.NewCanvas(320, 240)
	tests := []struct {
		lb1, lb2 *ui.Label
	}{
		// different colors
		{ui.NewLabel(c, "label", ui.WithColor(0x110011)), ui.NewLabel(c, "label", ui.WithColor(0x001100))},
		{ui.NewLabel(c, "label", ui.WithAlpha(0x01)), ui.NewLabel(c, "label", ui.WithAlpha(0x02))},
		// different font sizes
		{ui.NewLabel(c, "label", ui.WithFontSize(12)), ui.NewLabel(c, "label", ui.WithFontSize(12.5))},
		// different texts
		{ui.NewLabel(c, "label 1"), ui.NewLabel(c, "label 2")},
	}
	for _, test := range tests {
		require.NotEqual(t, test.lb1.Footprint(), test.lb2.Footprint())
	}
}
