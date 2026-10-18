//go:build regression

package ui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// https://github.com/gchumillas/prism/pull/77
func TestPull77(t *testing.T) {
	// https://github.com/gchumillas/prism/pull/77#discussion_r4302316204
	t.Run("footprint must include direction", func(t *testing.T) {
		ctx := NewContext(320, 240)
		hs := NewVStack(ctx)
		vs := NewHStack(ctx)
		require.NotEqual(t, hs.footprint(), vs.footprint())
	})
}
