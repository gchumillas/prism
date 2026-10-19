package ui

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAttachImageButton(t *testing.T) {
	c := NewCanvas(320, 240)
	vs := NewVStack(c)
	dat, err := os.ReadFile("testdata/icon1.png")
	require.NoError(t, err)
	_, err = AttachImageButtonTo(vs, bytes.NewReader(dat))
	require.NoError(t, err)
	require.Len(t, vs.children, 1)
}

func TestAttachLabelButton(t *testing.T) {
	c := NewCanvas(320, 240)
	vs := NewVStack(c)

	btn := NewLabelButton(c, "Click me!")
	_ = btn.AttachTo(vs)
	require.Len(t, vs.children, 1)
}
