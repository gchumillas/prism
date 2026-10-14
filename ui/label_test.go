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
		lb := ui.NewLabel(env.ctx, "Hello, World!")
		_, size := lb.Render(0, 0)
		require.Equal(t, image.Pt(181, 37), size)
	})

	t.Run("Contains", func(t *testing.T) {
		env := NewEnv()

		btn := ui.NewLabelButton(env.ctx, "Click me!")
		btn.Render(100, 50)

		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})

	t.Run("can be stacked vertically", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)

		btn := ui.NewLabelButton(env.ctx, "Click me!")
		vs.Add(btn)
		require.Contains(t, vs.Children(), btn)
	})
}
