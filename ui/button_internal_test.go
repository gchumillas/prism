package ui

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAttachImageButton(t *testing.T) {
	ctx := NewContext(320, 240)
	vs := NewVStack(ctx)
	dat, err := os.ReadFile("testdata/icon1.png")
	require.NoError(t, err)
	_, err = AttachImageButton(vs, bytes.NewReader(dat))
	require.NoError(t, err)
	require.Len(t, vs.children, 1)
}

func TestAttachLabelButton(t *testing.T) {
	ctx := NewContext(320, 240)
	vs := NewVStack(ctx)

	AttachLabelButton(vs, "Click me!")
	require.Len(t, vs.children, 1)
}
