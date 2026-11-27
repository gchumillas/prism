//go:build issues

package ui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

type UncomparableWidget []int // slices are not comparable

func (csw UncomparableWidget) Render(x, y int) []image.Rectangle { return nil }
func (csw UncomparableWidget) Size() image.Point                 { return image.Point{} }
func (csw UncomparableWidget) Footprint() []uint64               { return nil }

// https://github.com/gchumillas/prism/issues/146
func Test146(t *testing.T) {
	s := ui.NewHStack()
	csw := UncomparableWidget{}
	s.Add(csw)
	require.NotPanics(t, func() { s.Remove(csw) })
}

// https://github.com/gchumillas/prism/issues/158
func Test158(t *testing.T) {
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))
	s1 := ui.NewHStack()
	s2 := ui.NewHStack()
	btn := ui.NewButton(NewRectangle(c, 50, 50, color.Black))

	s1.Add(btn)
	s1.Render(0, 0)

	s2.Add(btn)
	s2.Render(100, 100)

	require.True(t, btn.Contains(125, 125))
	require.True(t, btn.Contains(25, 25))
}
