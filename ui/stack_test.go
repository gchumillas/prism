package ui_test

import (
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestStack_AttachTo(t *testing.T) {
	t.Run("self reference", func(t *testing.T) {
		s := ui.NewHStack()
		require.ErrorIs(t, s.AttachTo(s), ui.ErrCircularRef)
	})
	t.Run("circular reference 1", func(t *testing.T) {
		s1 := ui.NewHStack()
		s2 := ui.NewVStack()
		require.NoError(t, s2.AttachTo(s1))
		require.ErrorIs(t, s1.AttachTo(s2), ui.ErrCircularRef)
	})

	t.Run("circular reference 2", func(t *testing.T) {
		s1 := ui.NewHStack()
		s2 := ui.NewVStack()
		s3 := ui.NewHStack()
		require.NoError(t, s3.AttachTo(s2))
		require.NoError(t, s2.AttachTo(s1))
		require.ErrorIs(t, s1.AttachTo(s3), ui.ErrCircularRef)
	})
}

func TestVStack_Remove(t *testing.T) {
	t.Run("zero items", func(t *testing.T) {
		vs := ui.NewVStack()
		require.Equal(t, vs.Remove(nil), 0)
	})

	t.Run("one item", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack()
		r := NewRectangle(env.canvas, 100, 50, color.White)
		vs.Add(r)
		require.Equal(t, vs.Remove(r), 1)
	})

	t.Run("two items", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack()
		r := NewRectangle(env.canvas, 100, 50, color.White)
		vs.Add(r)
		vs.Add(r)
		require.Equal(t, vs.Remove(r), 2)
	})
}

func TestStack_Clear(t *testing.T) {
	env := NewEnv()
	s := ui.NewHStack()

	btn := ui.NewButton(NewRectangle(env.canvas, 50, 50, color.Black))
	_ = btn.AttachTo(s)

	s1 := ui.NewVStack()
	_ = s1.AttachTo(s)

	r := NewRectangle(env.canvas, 50, 50, color.Black)
	s.Add(r)

	s.Clear()
	require.ErrorIs(t, btn.DetachFrom(s), ui.ErrNotAttached)
	require.ErrorIs(t, s1.DetachFrom(s), ui.ErrNotAttached)
	newParent := ui.NewHStack()
	require.NoError(t, btn.AttachTo(newParent))
	require.NoError(t, s1.AttachTo(newParent))
	require.Zero(t, s.Remove(r))
}
