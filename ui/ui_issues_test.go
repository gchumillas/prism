//go:build issues

package ui

import (
	"image"
	"testing"

	"github.com/stretchr/testify/require"
)

type UncomparableWidget []int // slices are not comparable

func (csw UncomparableWidget) Render(x, y int) []image.Rectangle { return nil }
func (csw UncomparableWidget) Size() image.Point                 { return image.Point{} }
func (csw UncomparableWidget) Footprint() []uint64               { return nil }

// https://github.com/gchumillas/prism/issues/146
func Test146(t *testing.T) {
	s := NewHStack()
	csw := UncomparableWidget{}
	s.Add(csw)
	require.NotPanics(t, func() { s.Remove(csw) })
}
