package app

import "image/color"

func uint32ToColor(val uint32) color.Color {
	var color color.NRGBA
	color.R = uint8(val >> 24)
	color.G = uint8(val >> 16)
	color.B = uint8(val >> 8)
	color.A = uint8(val)
	return color
}
