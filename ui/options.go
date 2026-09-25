package ui

import "image"

type (
	SceneOpt  interface{ applyScene(*Scene) }
	VStackOpt interface{ applyVStack(*VStack) }
	LabelOpt  interface{ applyLabel(*Label) }
	TextOpt   interface{ applyText(*Text) }
)

type (
	paddingOpt  image.Point
	gapOpt      int
	fontSizeOpt float64
)

func Padding(width, height int) paddingOpt   { return paddingOpt(image.Point{X: width, Y: height}) }
func (val paddingOpt) applyVStack(x *VStack) { x.padding = image.Point(val) }

func Gap(val int) gapOpt                 { return gapOpt(val) }
func (val gapOpt) applyVStack(x *VStack) { x.gap = int(val) }

func FontSize(val float64) fontSizeOpt      { return fontSizeOpt(val) }
func (val fontSizeOpt) applyLabel(x *Label) { x.fontSize = float64(val) }
func (val fontSizeOpt) applyScene(x *Scene) { x.fontSize = float64(val) }
func (val fontSizeOpt) applyText(x *Text)   { x.fontSize = float64(val) }
