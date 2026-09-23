package ui

type (
	SceneOpt       interface{ applyScene(*Scene) }
	VStackOpt      interface{ applyVStack(*VStack) }
	LabelOpt       interface{ applyLabel(*Label) }
	LabelButtonOpt interface{ applyLabelButton(*LabelButton) }
	TextOpt        interface{ applyText(*Text) }
)

type (
	gapOpt      int
	fontSizeOpt float64
	dpiOpt      float64
)

func Gap(val int) gapOpt                 { return gapOpt(val) }
func (val gapOpt) applyVStack(x *VStack) { x.gap = int(val) }

func FontSize(val float64) fontSizeOpt      { return fontSizeOpt(val) }
func (val fontSizeOpt) applyLabel(x *Label) { x.fontSize = float64(val) }
func (val fontSizeOpt) applyScene(x *Scene) { x.fontSize = float64(val) }
func (val fontSizeOpt) applyText(x *Text)   { x.fontSize = float64(val) }

func DPI(val float64) dpiOpt           { return dpiOpt(val) }
func (val dpiOpt) applyScene(x *Scene) { x.dpi = float64(val) }
