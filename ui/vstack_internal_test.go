package ui

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVStack_Render(t *testing.T) {
	appWidth, appHeight := 320, 240

	t.Run("unused prev states are purged", func(t *testing.T) {
		ctx := NewContext(appWidth, appHeight)

		for n := range 3 {
			vs := NewVStack(ctx)
			var shapes []*LabelButton
			for i := range 3 {
				shape := NewLabelButton(ctx, fmt.Sprintf("Button %d", i))
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
