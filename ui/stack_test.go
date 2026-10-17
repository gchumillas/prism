package ui_test

import (
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestStack_Add(t *testing.T) {
	t.Run("custom widget", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		cw := NewCustomWidget(env.ctx, "custom ui")
		vs.Add(cw)
		require.Contains(t, vs.Children(), cw)
	})
}

func TestVStack_Remove(t *testing.T) {
	t.Run("zero items", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		require.Equal(t, vs.Remove(nil), 0)
	})

	t.Run("one item", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		r := NewRect(env.canvas, 0, 0, 100, 50, color.White)
		vs.Add(r)
		require.Equal(t, vs.Remove(r), 1)
	})

	t.Run("two items", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		r := NewRect(env.canvas, 0, 0, 100, 50, color.White)
		vs.Add(r)
		vs.Add(r)
		require.Equal(t, vs.Remove(r), 2)
	})
}

func TestVStack_RemoveButton(t *testing.T) {
	t.Run("zero button", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		require.False(t, vs.RemoveButton(nil))
	})

	t.Run("one button", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		btn := vs.AddButton(NewRect(env.canvas, 0, 0, 100, 50, color.White))
		require.True(t, vs.RemoveButton(btn))
	})
}
