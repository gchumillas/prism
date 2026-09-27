package ui

type (
	paddingXOpt int
	paddingYOpt int
	gapOpt      int
	fontSizeOpt float64
)

func Px(val int) paddingXOpt           { return paddingXOpt(val) } // horizontal padding
func Py(val int) paddingYOpt           { return paddingYOpt(val) } // vertical padding
func Gap(val int) gapOpt               { return gapOpt(val) }
func FontSize(val float64) fontSizeOpt { return fontSizeOpt(val) }
