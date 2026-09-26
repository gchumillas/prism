package ui

import "image"

type (
	paddingOpt  image.Point
	gapOpt      int
	fontSizeOpt float64
)

// shared options
func Padding(width, height int) paddingOpt { return paddingOpt(image.Point{X: width, Y: height}) }
func Gap(val int) gapOpt                   { return gapOpt(val) }
func FontSize(val float64) fontSizeOpt     { return fontSizeOpt(val) }
