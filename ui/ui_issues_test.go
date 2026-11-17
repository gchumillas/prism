//go:build issues

package ui

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

// https://github.com/gchumillas/prism/issues/130
func TestIssue130(t *testing.T) {
	c := NewCanvas(320, 240)
	hs := NewHStack()
	hs1 := NewHStack()
	r := NewRectangle(c, 50, 50, color.Black)
	hs1.Add(r)
	hs1.AttachTo(hs)
	hs.Render(0, 0)

	hs1.Remove(r)
	regs, _ := hs1.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(0, 0, 50, 50),
	}, regs)
}
