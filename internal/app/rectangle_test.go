package app_test

import (
	"image"
	"slices"
	"testing"

	"github.com/gchumillas/prism/internal/app"
	"github.com/stretchr/testify/require"
)

func TestRectangle_Size(t *testing.T) {
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))
	r := app.NewRectangle(c, 100, 50, app.WithFillColor(app.White))
	size := r.Size()
	require.Equal(t, image.Pt(100, 50), size)
}

func TestRectangle_Footprint(t *testing.T) {
	type shape = struct {
		width, height int
		fill, stroke  uint32
		lineWidth     int
	}
	tests := []struct {
		name   string
		shapes [2]shape
	}{
		{"change fill color", [2]shape{{100, 50, app.White, app.Black, 3}, {100, 50, app.Black, app.Black, 3}}},
		{"change stroke color", [2]shape{{100, 50, app.White, app.Black, 3}, {100, 50, app.White, app.White, 3}}},
		{"change line width", [2]shape{{100, 50, app.White, app.Black, 3}, {100, 50, app.White, app.Black, 5}}},
		{"change width", [2]shape{{100, 50, app.White, app.Black, 3}, {50, 50, app.White, app.Black, 3}}},
		{"change height", [2]shape{{100, 50, app.White, app.Black, 3}, {100, 100, app.White, app.Black, 3}}},
		{"change all", [2]shape{{50, 100, app.White, app.Black, 3}, {100, 50, app.Black, app.Black, 3}}},
	}

	c := image.NewRGBA(image.Rect(0, 0, 320, 240))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s0, s1 := test.shapes[0], test.shapes[1]
			r0 := app.NewRectangle(c, s0.width, s0.height, app.WithFillColor(s0.fill), app.WithStrokeColor(s0.stroke), app.WithLineWidth(s0.lineWidth))
			r1 := app.NewRectangle(c, s1.width, s1.height, app.WithFillColor(s1.fill), app.WithStrokeColor(s1.stroke), app.WithLineWidth(s1.lineWidth))
			require.False(t, slices.Equal(r0.Footprint(), r1.Footprint()))
		})
	}
}
