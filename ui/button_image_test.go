package ui_test

import (
	"bytes"
	"image"
	"os"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestImageButton(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		dat, err := os.ReadFile("testdata/icon1.png")
		require.NoError(t, err)
		_, err = ui.AttachImageButtonTo(vs, env.ctx, bytes.NewReader(dat))
		require.NoError(t, err)
		_, size := vs.Render(0, 0)
		require.Equal(t, image.Pt(34, 34), size)
	})

	t.Run("Contains", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		dat, err := os.ReadFile("testdata/icon1.png")
		require.NoError(t, err)
		btn, err := ui.AttachImageButtonTo(vs, env.ctx, bytes.NewReader(dat))
		require.NoError(t, err)
		vs.Render(100, 50)

		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})
}
