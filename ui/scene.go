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

func (s *Scene) SetDefaultFontSize(val float64) { s.fontSize = val }
func (s *Scene) DefaultFontSize() float64       { return s.fontSize }

func (s *Scene) SetDPI(val float64) { s.dpi = val }
func (s *Scene) DPI() float64       { return s.dpi }
