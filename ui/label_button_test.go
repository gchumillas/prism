package ui_test

import (
	"image"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestLabelButton_Render(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		env := NewEnv()
		btn := ui.NewLabelButton(env.ctx, "Hello, World!")
		_, size := btn.Render(0, 0)
		require.Equal(t, image.Pt(257, 75), size)
	})
}
