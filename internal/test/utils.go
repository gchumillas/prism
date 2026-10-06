package test

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDir = "testdata"

func TakeSnapshot(t testing.TB, img image.Image) {
	name := filepath.Join(testDir, t.Name()) + ".png"
	dir := strings.TrimSuffix(name, filepath.Base(name))
	err := os.MkdirAll(dir, 0o755)
	require.NoError(t, err)
	f, err := os.Create(name)
	require.NoError(t, err)
	defer func() { _ = f.Close() }()
	err = png.Encode(f, img)
	require.NoError(t, err)
}

func rectIsValid(t testing.TB, x0, y0, x1, y1 int) bool {
	t.Helper()
	if x0 > x1 || y0 > y1 {
		t.Errorf("Invalid rectangle")
		return false
	}
	return true
}

func RectIsColor(t testing.TB, img image.Image, x0, y0, x1, y1 int, c color.Color) bool {
	t.Helper()
	if !rectIsValid(t, x0, y0, x1, y1) {
		return false
	}
	c = img.ColorModel().Convert(c)
	b := img.Bounds()
	x0, y0 = max(x0, b.Min.X), max(y0, b.Min.Y)
	x1, y1 = min(x1, b.Max.X-1), min(y1, b.Max.Y-1)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if d := img.At(x, y); d != c {
				t.Errorf("Want %v, got %v", c, d)
				return false
			}
		}
	}
	return true
}

func RectHasSomeColors(t testing.TB, img image.Image, x0, y0, x1, y1 int) {
	t.Helper()
	if !rectIsValid(t, x0, y0, x1, y1) {
		return
	}
	b := img.Bounds()
	x0, y0 = max(x0, b.Min.X), max(y0, b.Min.Y)
	x1, y1 = min(x1, b.Max.X-1), min(y1, b.Max.Y-1)
	curColor := img.At(x0, y0)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if img.At(x, y) != curColor {
				return
			}
		}
	}
	t.Error("Want some colors, got a single color")
}

func HlineIsColor(t testing.TB, img image.Image, x0, y0, width int, c color.Color) {
	t.Helper()
	RectIsColor(t, img, x0, y0, x0+width-1, y0, c)
}

func VlineIsColor(t testing.TB, img image.Image, x0, y0, height int, c color.Color) {
	t.Helper()
	RectIsColor(t, img, x0, y0, x0, y0+height-1, c)
}

func RectBorderIsColor(t testing.TB, img image.Image, x0, y0, x1, y1 int, c color.Color) {
	t.Helper()
	HlineIsColor(t, img, x0, y0, x1-x0, c)
	HlineIsColor(t, img, x0, y1, x1-x0, c)
	VlineIsColor(t, img, x0, y0, y1-y0, c)
	VlineIsColor(t, img, x0, y1, y1-y0, c)
}

func DotHasSomeColors(t testing.TB, img image.Image, x, y, size int) {
	t.Helper()
	RectHasSomeColors(t, img, x-size, y-size, x+size, y+size)
}
