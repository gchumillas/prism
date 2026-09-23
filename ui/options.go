package ui

type (
	SceneOpt       interface{ applyScene(*Scene) }
	VStackOpt      interface{ applyVStack(*VStack) }
	LabelOpt       interface{ applyLabel(*Label) }
	LabelButtonOpt interface{ applyLabelButton(*LabelButton) }
	TextOpt        interface{ applyText(*Text) }
)

type (
	xOpt        int
	yOpt        int
	widthOpt    int
	gapOpt      int
	fontSizeOpt float64
	dpiOpt      float64
)

func X(val int) xOpt                             { return xOpt(val) }
func (val xOpt) applyLabel(x *Label)             { x.position.X = int(val) }
func (val xOpt) applyLabelButton(x *LabelButton) { x.position.Y = int(val) }

func Y(val int) yOpt                             { return yOpt(val) }
func (val yOpt) applyLabel(x *Label)             { x.position.Y = int(val) }
func (val yOpt) applyLabelButton(x *LabelButton) { x.position.Y = int(val) }

func Width(val int) widthOpt                         { return widthOpt(val) }
func (val widthOpt) applyLabelButton(x *LabelButton) { x.width = int(val) }

func Gap(val int) gapOpt                 { return gapOpt(val) }
func (val gapOpt) applyVStack(x *VStack) { x.gap = int(val) }

func FontSize(val float64) fontSizeOpt      { return fontSizeOpt(val) }
func (val fontSizeOpt) applyLabel(x *Label) { x.fontSize = float64(val) }
func (val fontSizeOpt) applyScene(x *Scene) { x.fontSize = float64(val) }
func (val fontSizeOpt) applyText(x *Text)   { x.fontSize = float64(val) }

func DPI(val float64) dpiOpt           { return dpiOpt(val) }
func (val dpiOpt) applyScene(x *Scene) { x.dpi = float64(val) }
