package ui_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestImageButton(t *testing.T) {
	t.Run("Contains", func(t *testing.T) {
		env := NewEnv()
		dat, err := os.ReadFile("testdata/icon1.png")
		require.NoError(t, err)
		btn, err := ui.LoadImageButton(env.ctx, bytes.NewReader(dat))
		require.NoError(t, err)
		btn.Render(100, 50)

		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})

	t.Run("can be stacked vertically", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)

		dat, err := os.ReadFile("testdata/icon1.png")
		require.NoError(t, err)
		btn, err := ui.LoadImageButton(env.ctx, bytes.NewReader(dat))
		require.NoError(t, err)
		vs.Add(btn)
		require.Contains(t, vs.Children(), btn)
	})
}
