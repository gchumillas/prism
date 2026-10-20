package ui_test

import (
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestButton_AttachTo(t *testing.T) {
	env := NewEnv()
	hs := ui.NewHStack(env.canvas)
	btn := ui.NewButton(ui.NewRectangle(env.canvas, 100, 50, color.White))
	require.NoError(t, btn.AttachTo(hs))
	require.ErrorIs(t, btn.AttachTo(hs), ui.ErrAttached)
}

func TestButton_DetachFrom(t *testing.T) {
	t.Run("zero button", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		hs := ui.NewHStack(env.canvas)
		btn := ui.NewButton(ui.NewRectangle(env.canvas, 100, 50, color.White))
		_ = btn.AttachTo(vs)
		require.ErrorIs(t, ui.ErrNotAttached, btn.DetachFrom(hs))
	})

	t.Run("one button", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		btn := ui.NewButton(NewRect(env.image, 0, 0, 100, 50, color.White))
		_ = btn.AttachTo(vs)
		require.NoError(t, btn.DetachFrom(vs))
		require.ErrorIs(t, ui.ErrNotAttached, btn.DetachFrom(vs))
	})
}
