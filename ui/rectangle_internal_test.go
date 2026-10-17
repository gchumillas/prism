package ui

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRectangle(t *testing.T) {
	t.Run("can be stacked vertically", func(t *testing.T) {
		ctx := NewContext(320, 240)
		vs := NewVStack(ctx)
		r := NewRectangle(ctx, 50, 100, color.White)
		vs.Add(r)
		require.Contains(t, vs.children, r)
	})
}
