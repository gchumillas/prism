package ui_test

import (
	"bytes"
	"image"
	"os"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestImage_Render(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		env := NewEnv()
		icon, err := os.ReadFile("testdata/icon1.png")
		require.NoError(t, err)
		img, err := ui.LoadImage(env.canvas, bytes.NewReader(icon))
		require.NoError(t, err)
		_, size := img.Render(0, 0)
		require.Equal(t, image.Pt(24, 24), size)
	})
}
