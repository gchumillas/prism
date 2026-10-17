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
		vs := ui.NewVStack(env.ctx)
		ui.AttachLabelButton(vs, env.ctx, "Hello, World!")
		_, size := vs.Render(0, 0)
		require.Equal(t, image.Pt(257, 75), size)
	})

	t.Run("Contains", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		btn := ui.AttachLabelButton(vs, env.ctx, "Click me!")
		vs.Render(100, 50)

		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})
}
