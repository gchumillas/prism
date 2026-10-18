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
		vs := ui.NewVStack(env.canvas)
		hs := ui.NewHStack(env.canvas)
		btn := ui.AttachLabelButtonTo(vs, env.canvas, "Click me!")
		require.False(t, btn.DetachFrom(hs))
	})

	t.Run("one button", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		btn := ui.AttachButtonTo(vs, NewRect(env.image, 0, 0, 100, 50, color.White))
		require.True(t, btn.DetachFrom(vs))
	})
}
