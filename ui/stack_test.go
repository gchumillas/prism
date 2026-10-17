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
	env := NewEnv()
	vs := ui.NewVStack(env.ctx)
	r := NewRect(env.canvas, 0, 0, 100, 50, color.White)
	vs.Add(r)
	require.Contains(t, vs.Children(), r)
	vs.Remove(r)
	require.NotContains(t, vs.Children(), r)
}

func TestVStack_RemoveButton(t *testing.T) {
	env := NewEnv()
	vs := ui.NewVStack(env.ctx)
	btn := vs.AddButton(NewRect(env.canvas, 0, 0, 100, 50, color.White))
	require.Len(t, vs.Children(), 1)
	vs.RemoveButton(btn)
	require.Empty(t, vs.Children())
}
