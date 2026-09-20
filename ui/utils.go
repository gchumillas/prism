package ui

import (
	"image/color"
	"image/draw"
)

func DrawHLine(img draw.Image, x0, y0, width int, c color.Color) {
	for w := range width {
		img.Set(x0+w, y0, c)
	}
}

func DrawVLine(img draw.Image, x0, y0, height int, c color.Color) {
	for h := range height {
		img.Set(x0, y0+h, c)
	}
}

func DrawRect(img draw.Image, x0, y0, x1, y1 int, c color.Color) {
	w, h := x1-x0, y1-y0
	DrawHLine(img, x0, y0, w, c) // top line
	DrawHLine(img, x0, y1, w, c) // bottom line
	DrawVLine(img, x0, y0, h, c) // left line
	DrawVLine(img, x1, y0, h, c) // right line
}
