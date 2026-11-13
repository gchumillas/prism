package ui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestInset_Render(t *testing.T) {
	env := NewEnv()
	r := ui.NewRectangle(env.canvas, 30, 25, color.Black)
	i := ui.NewInset(r, 5, 10)
	_, size := i.Render(0, 0)
	require.Equal(t, image.Pt(40, 45), size)
}
