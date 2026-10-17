//go:build regression

package ui_test

import (
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

// https://github.com/gchumillas/prism/pull/77
func TestPull77(t *testing.T) {
	// https://github.com/gchumillas/prism/pull/77#discussion_r4302316204
	t.Run("footprint must include direction", func(t *testing.T) {
		ctx := NewEnv()
		hs := ui.NewVStack(ctx.ctx)
		vs := ui.NewHStack(ctx.ctx)
		require.NotEqual(t, hs.Footprint(), vs.Footprint())
	})
}
