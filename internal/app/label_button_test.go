package app_test

import (
	"image"
	"testing"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

type Env struct {
	image  *image.RGBA
	canvas *ui.Canvas
}

func NewEnv() *Env {
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	c := ui.NewCanvasFromRGBA(img)
	return &Env{image: img, canvas: c}
}

func TestLabelButton_Render(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		btn := app.NewLabelButton(env.canvas, "Hello, World!")
		_ = btn.AttachTo(vs)
		_, size := vs.Render(0, 0)
		require.Equal(t, image.Pt(257, 75), size)
	})

	t.Run("Contains", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.canvas)
		btn := app.NewLabelButton(env.canvas, "Click me!")
		_ = btn.AttachTo(vs)
		vs.Render(100, 50)

		require.False(t, btn.Contains(0, 0))   // outside
		require.True(t, btn.Contains(100, 50)) // inside (left-top corner)
		require.True(t, btn.Contains(110, 60)) // inside
	})
}
