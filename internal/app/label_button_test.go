package app_test

import (
	"image"
	"testing"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestLabelButton(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		vs := ui.NewVStack()
		btn := app.NewLabelButton(c, "Hello, World!", 10)
		_ = btn.AttachTo(vs)
		size := vs.Size()
		require.Equal(t, image.Pt(257, 75), size)
	})

	t.Run("Contains", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		vs := ui.NewVStack()
		btn := app.NewLabelButton(c, "Click me!", 10)
		_ = btn.AttachTo(vs)
		vs.Render(100, 50)

		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})
}
