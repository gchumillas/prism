package ui_test

import (
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestStack_AttachTo(t *testing.T) {
	t.Run("self reference", func(t *testing.T) {
		env := NewEnv()
		s := ui.NewHStack(env.canvas)
		require.ErrorIs(t, s.AttachTo(s), ui.ErrCircularRef)
	})
	t.Run("circular reference 1", func(t *testing.T) {
		env := NewEnv()
		s1 := ui.NewHStack(env.canvas)
		s2 := ui.NewVStack(env.canvas)
		require.NoError(t, s2.AttachTo(s1))
		require.ErrorIs(t, s1.AttachTo(s2), ui.ErrCircularRef)
	})

	t.Run("circular reference 2", func(t *testing.T) {
		env := NewEnv()
		s1 := ui.NewHStack(env.canvas)
		s2 := ui.NewVStack(env.canvas)
		s3 := ui.NewHStack(env.canvas)
		require.NoError(t, s3.AttachTo(s2))
		require.NoError(t, s2.AttachTo(s1))
		require.ErrorIs(t, s1.AttachTo(s3), ui.ErrCircularRef)
	})
}

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
