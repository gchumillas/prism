package ui_test

import (
	"image"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestLabelButton(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		env := NewEnv()
		lb := ui.NewLabel(env.canvas, "Hello, World!")
		_, size := lb.Render(0, 0)
		require.Equal(t, image.Pt(181, 37), size)
	})
}
