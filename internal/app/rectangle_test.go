package app_test

import (
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/gchumillas/prism/internal/app"
	"github.com/stretchr/testify/require"
)

func TestRectangle_Render(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		r := app.NewRectangle(c, 100, 50, color.White)
		size := r.Size()
		require.Equal(t, image.Pt(100, 50), size)
	})
}

func TestRectangle_Footprint(t *testing.T) {
	type shape = struct {
		width, height int
		color         color.Color
	}
	tests := []struct {
		name   string
		shapes [2]shape
	}{
		{"change color", [2]shape{{100, 50, color.White}, {100, 50, color.Black}}},
		{"change width", [2]shape{{100, 50, color.White}, {100, 50, color.Black}}},
		{"change height", [2]shape{{100, 50, color.White}, {100, 50, color.Black}}},
		{"change all", [2]shape{{100, 50, color.White}, {100, 50, color.Black}}},
	}

	c := image.NewRGBA(image.Rect(0, 0, 320, 240))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s0, s1 := test.shapes[0], test.shapes[1]
			r0 := app.NewRectangle(c, s0.width, s0.height, s0.color)
			r1 := app.NewRectangle(c, s1.width, s1.height, s1.color)
			require.False(t, slices.Equal(r0.Footprint(), r1.Footprint()))
		})
	}
}
