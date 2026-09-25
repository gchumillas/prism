package test

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/fogleman/gg"
	"github.com/stretchr/testify/require"
)

const testDir = "testdata"

func GenSnapshot(t *testing.T, canvas *gg.Context) *gg.Context {
	if os.Getenv("GENERATE_SNAPSHOTS") == "1" {
		name := filepath.Join(testDir, t.Name()+".png")
		err := os.MkdirAll(testDir, 0o755)
		require.NoError(t, err)
		f, err := os.Create(name)
		require.NoError(t, err)
		defer func() { _ = f.Close() }()
		err = canvas.EncodePNG(f)
		require.NoError(t, err)
	}
	return canvas
}

func rectIsValid(t *testing.T, x0, y0, x1, y1 int) bool {
	t.Helper()
	if x0 > x1 || y0 > y1 {
		t.Errorf("Invalid rectangle")
		return false
	}
	return true
}

func rectIsColor(t *testing.T, canvas *gg.Context, x0, y0, x1, y1 int, c color.Color) bool {
	t.Helper()
	if !rectIsValid(t, x0, y0, x1, y1) {
		return false
	}
	img := canvas.Image()
	c = img.ColorModel().Convert(c)
	b := img.Bounds()
	x0, y0 = max(x0, b.Min.X), max(y0, b.Min.Y)
	x1, y1 = min(x1, b.Max.X-1), min(y1, b.Max.Y-1)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if d := img.At(x, y); d != c {
				t.Error("Want single color, got multiple")
				return false
			}
		}
	}
	return true
}

func RectHasSomeColors(t *testing.T, canvas *gg.Context, x0, y0, x1, y1 int) {
	t.Helper()
	if !rectIsValid(t, x0, y0, x1, y1) {
		return
	}

	img := canvas.Image()
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
	t.Error("Want single color, got multiple")
}

func HlineIsColor(t *testing.T, canvas *gg.Context, x0, y0, width int, c color.Color) {
	t.Helper()
	x1, y1 := x0+width-1, y0
	rectIsColor(t, canvas, x0, y0, x1, y1, c)
}

func VlineIsColor(t *testing.T, canvas *gg.Context, x0, y0, height int, c color.Color) {
	t.Helper()
	x1, y1 := x0, y0+height-1
	rectIsColor(t, canvas, x0, y0, x1, y1, c)
}

func DotHasSomeColors(t *testing.T, canvas *gg.Context, x, y, size int) {
	t.Helper()
	x0, y0, x1, y1 := x-size, y-size, x+size, y+size
	RectHasSomeColors(t, canvas, x0, y0, x1, y1)
}
