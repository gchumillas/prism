package ui_test

import (
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestText_Footprint(t *testing.T) {
	c := ui.NewCanvas(320, 240)
	tests := []struct {
		lb1, lb2 *ui.Text
	}{
		// different colors
		{ui.NewText(c, "lorem ipsum", 200, ui.WithColor(0x110011)), ui.NewText(c, "lorem ipsum", 200, ui.WithColor(0x001100))},
		{ui.NewText(c, "lorem ipsum", 200, ui.WithAlpha(0x01)), ui.NewText(c, "lorem ipsum", 200, ui.WithAlpha(0x02))},
		// different font sizes
		{ui.NewText(c, "lorem ipsum", 200, ui.WithFontSize(12)), ui.NewText(c, "lorem ipsum", 200, ui.WithFontSize(12.5))},
		// different texts
		{ui.NewText(c, "lorem ipsum", 200), ui.NewText(c, "lorem ipsum dolor", 200)},
		// different widths
		{ui.NewText(c, "lorem ipsum", 200), ui.NewText(c, "lorem ipsum", 300)},
	}
	for _, test := range tests {
		require.NotEqual(t, test.lb1.Footprint(), test.lb2.Footprint())
	}
}
