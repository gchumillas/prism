package ui_test

import (
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestRectangle(t *testing.T) {
	t.Run("can be stacked vertically", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)

		r := ui.NewRectangle(env.ctx, 50, 100, color.White)
		vs.Add(r)
		require.Contains(t, vs.Children(), r)
	})
}

func TestRectangle_Render(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		env := NewEnv()
		r := ui.NewRectangle(env.ctx, 100, 50, color.White)
		_, size := r.Render(0, 0)
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

	env := NewEnv()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s0, s1 := test.shapes[0], test.shapes[1]
			r0 := ui.NewRectangle(env.ctx, s0.width, s0.height, s0.color)
			r1 := ui.NewRectangle(env.ctx, s1.width, s1.height, s1.color)
			require.False(t, slices.Equal(r0.Footprint(), r1.Footprint()))
		})
	}
}
