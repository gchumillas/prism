package ui_test

import (
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestVStack_Remove(t *testing.T) {
	t.Run("zero items", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		require.Equal(t, vs.Remove(nil), 0)
	})

	t.Run("one item", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		r := NewRect(env.image, 0, 0, 100, 50, color.White)
		vs.Add(r)
		require.Equal(t, vs.Remove(r), 1)
	})

	t.Run("two items", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		r := NewRect(env.image, 0, 0, 100, 50, color.White)
		vs.Add(r)
		vs.Add(r)
		require.Equal(t, vs.Remove(r), 2)
	})
}
