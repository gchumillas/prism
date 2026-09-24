package ui

type (
	SceneOpt  interface{ applyScene(*Scene) }
	VStackOpt interface{ applyVStack(*VStack) }
	LabelOpt  interface{ applyLabel(*Label) }
	TextOpt   interface{ applyText(*Text) }
)

type (
	gapOpt      float64
	fontSizeOpt float64
)

func Gap(val int) gapOpt                 { return gapOpt(val) }
func (val gapOpt) applyVStack(x *VStack) { x.gap = float64(val) }

func FontSize(val float64) fontSizeOpt      { return fontSizeOpt(val) }
func (val fontSizeOpt) applyLabel(x *Label) { x.fontSize = float64(val) }
func (val fontSizeOpt) applyScene(x *Scene) { x.fontSize = float64(val) }
func (val fontSizeOpt) applyText(x *Text)   { x.fontSize = float64(val) }
