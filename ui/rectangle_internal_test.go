package ui

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRectangle(t *testing.T) {
	t.Run("can be stacked vertically", func(t *testing.T) {
		c := NewCanvas(320, 240)
		vs := NewVStack(c)
		r := NewRectangle(c, 50, 100, color.White)
		vs.Add(r)
		require.Contains(t, vs.children, r)
	})
}
