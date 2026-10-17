package ui

import (
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

type CustomWidget struct {
	BaseWidget
	Label Label
}

func NewCustomWidget(c Context, text string) *CustomWidget {
	return &CustomWidget{Label: *NewLabel(c, text)}
}

func (cw *CustomWidget) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	return cw.Label.Render(x, y)
}

func (cw *CustomWidget) Footprint() []uint64 {
	return cw.Label.Footprint()
}

func TestStack_Add(t *testing.T) {
	t.Run("custom widget", func(t *testing.T) {
		ctx := NewContext(320, 240)
		vs := NewVStack(ctx)
		cw := NewCustomWidget(ctx, "custom ui")
		vs.Add(cw)
		require.Contains(t, vs.children, cw)
	})
}

func TestVStack_Render(t *testing.T) {
	appWidth, appHeight := 320, 240

	t.Run("unused prev states are purged", func(t *testing.T) {
		ctx := NewContext(appWidth, appHeight)

		for n := range 3 {
			vs := NewVStack(ctx)
			var shapes []*Rectangle
			for i := range 3 {
				shape := NewRectangle(ctx, 100, 50, color.Gray{Y: uint8(i)})
				shapes = append(shapes, shape)
				vs.Add(shape)
			}
			vs.Render(0, 0)
			require.Len(t, vs.prevChildStates, 3)

			vs.Remove(shapes[n])
			vs.Render(0, 0)
			require.Len(t, vs.prevChildStates, 2)
			// all shapes, except the removed one, are in prevChildStates
			for i := range shapes {
				found := slices.ContainsFunc(vs.prevChildStates, func(prev childState) bool {
					return slices.Equal(prev.footprint, shapes[i].Footprint())
				})
				if i == n {
					require.False(t, found, "removed shape should be purged")
				} else {
					require.True(t, found, "shape %d should be kept", i)
				}
			}
		}
	})
}
