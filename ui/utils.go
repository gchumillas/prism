package ui

import (
	"image/color"
	"image/draw"
)

func drawHLine(img draw.Image, x0, y0, width int, c color.Color) {
	for w := range width {
		img.Set(x0+w, y0, c)
	}
}

func drawVLine(img draw.Image, x0, y0, height int, c color.Color) {
	for h := range height {
		img.Set(x0, y0+h, c)
	}
}

func drawRect(img draw.Image, x0, y0, x1, y1 int, c color.Color) {
	w, h := x1-x0, y1-y0
	drawHLine(img, x0, y0, w, c) // top line
	drawHLine(img, x0, y1, w, c) // bottom line
	drawVLine(img, x0, y0, h, c) // left line
	drawVLine(img, x1, y0, h, c) // right line
}
