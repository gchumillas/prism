package ui

import "image/draw"

const (
	defaultFontSize = 10
	defaultDPI      = 226 // 96: Web; 226: reMarkable 2
)

type Scene struct {
	canvas   draw.Image
	fontSize float64
	dpi      float64
}

func NewScene(canvas draw.Image) *Scene {
	return &Scene{
		canvas:   canvas,
		fontSize: defaultFontSize,
		dpi:      defaultDPI,
	}
}

func (sc *Scene) SetDefaultFontSize(val float64) { sc.fontSize = val }
func (sc *Scene) DefaultFontSize() float64       { return sc.fontSize }

func (sc *Scene) SetDPI(val float64) { sc.dpi = val }
func (sc *Scene) DPI() float64       { return sc.dpi }
