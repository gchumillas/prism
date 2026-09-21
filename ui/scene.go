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

func NewScene(canvas draw.Image, opts ...SceneOpt) *Scene {
	s := &Scene{
		canvas:   canvas,
		fontSize: defaultFontSize,
		dpi:      defaultDPI,
	}
	for _, opt := range opts {
		opt.applyScene(s)
	}
	return s
}
