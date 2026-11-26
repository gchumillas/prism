package ui_test

import (
	"bytes"
	"image"
	"os"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

func TestImage_LoadImageFromPNG(t *testing.T) {
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))
	icon, err := os.ReadFile("testdata/icon.png")
	require.NoError(t, err)
	img, err := ui.LoadImageFromPNG(c, bytes.NewReader(icon))
	require.NoError(t, err)
	size := img.Size()
	require.Equal(t, image.Pt(24, 24), size)
}

func TestImage_LoadImageFromJPEG(t *testing.T) {
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))
	icon, err := os.ReadFile("testdata/icon.jpg")
	require.NoError(t, err)
	img, err := ui.LoadImageFromJPEG(c, bytes.NewReader(icon))
	require.NoError(t, err)
	size := img.Size()
	require.Equal(t, image.Pt(24, 24), size)
}
