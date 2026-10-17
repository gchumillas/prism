package ui_test

import (
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestButton_DetachFrom(t *testing.T) {
	t.Run("zero button", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		hs := ui.NewHStack(env.ctx)
		btn := ui.AttachLabelButton(vs, env.ctx, "Click me!")
		require.False(t, btn.DetachFrom(hs))
	})

	t.Run("one button", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		btn := ui.AttachButton(vs, NewRect(env.canvas, 0, 0, 100, 50, color.White))
		require.True(t, btn.DetachFrom(vs))
	})
}
