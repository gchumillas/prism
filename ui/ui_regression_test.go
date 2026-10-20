//go:build regression

package ui

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/gchumillas/prism/internal/test"
	"github.com/stretchr/testify/require"
)

// https://github.com/gchumillas/prism/pull/77
func TestPull77(t *testing.T) {
	// https://github.com/gchumillas/prism/pull/77#discussion_r4302316204
	t.Run("footprint must include direction", func(t *testing.T) {
		c := NewCanvas(320, 240)
		hs := NewVStack(c)
		vs := NewHStack(c)
		require.NotEqual(t, hs.footprint(), vs.footprint())
	})
}

// https://github.com/gchumillas/prism/issues/110
func TestIssue110(t *testing.T) {
	width := 50
	x0, y0 := 100, 50
	x1, y1 := x0+50, y0+50
	c := NewCanvas(320, 240)

	// internal stack
	hs1 := NewHStack(c)
	r := NewRectangle(c, width, width, color.White)
	hs1.Add(r)

	// render a external stack that contains a internal stack
	hs := NewHStack(c)
	_ = hs1.AttachTo(hs)
	hs.Render(x0, y0)

	// remove the internal stack, render and clear "dirty regions"
	_ = hs1.DetachFrom(hs)
	dirtyRegs, _ := hs.Render(x0, y0)
	for _, reg := range dirtyRegs {
		draw.Draw(c.image, reg, image.NewUniform(color.Transparent), image.Point{}, draw.Src)
	}
	test.RectIsColor(t, c.image, x0, y0, x1, y1, color.Transparent) // has been cleared correctly

	// add the internal stack again and render
	_ = hs1.AttachTo(hs)
	hs.Render(x0, y0)
	test.RectIsColor(t, c.image, x0+10, y0+10, x1-10, y1-10, color.White) // has been rendered correctly

	// remove the internal stack again, repaint and clear "dirty regions"
	hs1.DetachFrom(hs)
	dirtyRegs, _ = hs.Render(x0, y0)
	for _, reg := range dirtyRegs {
		draw.Draw(c.image, reg, image.NewUniform(color.Transparent), image.Point{}, draw.Src)
	}
	test.RectIsColor(t, c.image, x0, y0, x1, y1, color.Transparent) // has been cleared correctly
}
