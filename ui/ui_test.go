package ui_test

import (
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"

	_ "image/png"
)

const (
	appWidth  = 640
	appHeight = 480
)

type App struct {
	*ui.VStack
	bg *ui.Rectangle
}

type (
	AppConfig struct {
		padding image.Point
		gap     int
	}

	AppOption func(*AppConfig)
)

func WithPadding(val image.Point) AppOption {
	return func(cfg *AppConfig) {
		cfg.padding = val
	}
}

func WithGap(val int) AppOption {
	return func(cfg *AppConfig) {
		cfg.gap = val
	}
}

func NewApp(ctx ui.Context, w, h int, opts ...AppOption) *App {
	cfg := &AppConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return &App{
		VStack: ui.NewVStack(ctx, ui.WithPadding(cfg.padding), ui.WithGap(cfg.gap)),
		bg:     ui.NewRectangle(ctx, w, h, color.White),
	}
}

func (s *App) Render() {
	s.bg.Render(0, 0)
	s.VStack.Render(0, 0)
}

func TestMergeRegions(t *testing.T) {
	mirrorHor := func(rs []image.Rectangle) {
		for i, r := range rs {
			rs[i] = image.Rectangle{
				Min: image.Pt(-r.Min.X, r.Min.Y),
				Max: image.Pt(-r.Max.X, r.Max.Y),
			}.Canon()
		}
	}

	mirrorVer := func(rs []image.Rectangle) {
		for i, r := range rs {
			rs[i] = image.Rectangle{
				Min: image.Pt(r.Min.X, -r.Min.Y),
				Max: image.Pt(r.Max.X, -r.Max.Y),
			}.Canon()
		}
	}

	movements := []func([]image.Rectangle){func(r []image.Rectangle) {}, mirrorHor, mirrorVer, mirrorHor}

	t.Run("disjoint rectangles", func(t *testing.T) {
		tests := [][]image.Rectangle{
			{image.Rect(-2, -2, 2, 2), image.Rect(3, 3, 5, 5)},
			{image.Rect(-2, -2, 2, 2), image.Rect(3, 1, 5, 3)},
			{image.Rect(-2, -2, 2, 2), image.Rect(3, -1, 5, 1)},
			{image.Rect(-2, -2, 2, 2), image.Rect(3, -3, 5, 3)},
		}
		for _, test := range tests {
			for _, m := range movements {
				m(test)
				output := ui.MergeRegions(test)
				require.ElementsMatch(t, test, output)
			}
		}
	})

	t.Run("intersecting rectangles", func(t *testing.T) {
		tests := []struct {
			input    []image.Rectangle
			expected []image.Rectangle
		}{
			{
				input:    []image.Rectangle{image.Rect(-2, -2, 2, 2), image.Rect(1, 1, 3, 3)},
				expected: []image.Rectangle{image.Rect(-2, -2, 2, 1), image.Rect(-2, 1, 3, 2), image.Rect(1, 2, 3, 3)},
			},
			{
				input:    []image.Rectangle{image.Rect(-2, -2, 2, 2), image.Rect(-1, -1, 3, 1)},
				expected: []image.Rectangle{image.Rect(-2, -2, 2, -1), image.Rect(-2, -1, 3, 1), image.Rect(-2, 1, 2, 2)},
			},
		}
		for _, test := range tests {
			for _, m := range movements {
				m(test.input)
				m(test.expected)
				output := ui.MergeRegions(test.input)
				require.ElementsMatch(t, test.expected, output)
			}
		}
	})

	t.Run("nested rectangles", func(t *testing.T) {
		input := []image.Rectangle{image.Rect(-2, -2, 2, 2), image.Rect(-1, -1, 1, 1)}
		expected := []image.Rectangle{image.Rect(-2, -2, 2, 2)}

		for range 2 {
			output := ui.MergeRegions(input)
			require.Equal(t, expected, output)
			slices.Reverse(input)
			slices.Reverse(expected)
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		tests := []struct {
			input    []image.Rectangle
			expected []image.Rectangle
		}{
			{
				input:    []image.Rectangle{image.Rect(-1, -1, 1, 1), image.Rect(1, -1, 3, 1)},
				expected: []image.Rectangle{image.Rect(-1, -1, 3, 1)},
			},
			{
				input:    []image.Rectangle{image.Rect(-1, -1, 1, 1), image.Rect(1, 1, 3, 3)},
				expected: []image.Rectangle{image.Rect(-1, -1, 1, 1), image.Rect(1, 1, 3, 3)},
			},
		}
		for _, test := range tests {
			output := ui.MergeRegions(test.input)
			require.ElementsMatch(t, test.expected, output)
		}
	})

	t.Run("combined", func(t *testing.T) {
		input := []image.Rectangle{
			image.Rect(-2, -2, 2, 2),
			image.Rect(1, 0, 4, 4),
			image.Rect(2, 1, 3, 2),
			image.Rect(3, 3, 5, 5),
			image.Rect(0, 4, 1, 5),
		}
		expected := []image.Rectangle{
			image.Rect(-2, -2, 2, 0),
			image.Rect(-2, 0, 4, 2),
			image.Rect(1, 2, 4, 3),
			image.Rect(1, 3, 5, 4),
			image.Rect(0, 4, 1, 5),
			image.Rect(3, 4, 5, 5),
		}

		output := ui.MergeRegions(input)
		require.ElementsMatch(t, expected, output)
	})
}
