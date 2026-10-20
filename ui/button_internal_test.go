package ui

import (
	"bytes"
	"image/color"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAttachImageButton(t *testing.T) {
	c := NewCanvas(320, 240)
	vs := NewVStack(c)
	dat, err := os.ReadFile("testdata/icon1.png")
	require.NoError(t, err)
	btn, err := LoadImageButton(c, bytes.NewReader(dat))
	require.NoError(t, err)
	_ = btn.AttachTo(vs)
	require.Len(t, vs.children, 1)
}

func TestAttachLabelButton(t *testing.T) {
	c := NewCanvas(320, 240)
	vs := NewVStack(c)

	btn := NewButton(NewRectangle(c, 100, 50, color.White))
	_ = btn.AttachTo(vs)
	require.Len(t, vs.children, 1)
}
