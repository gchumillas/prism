package ui

import (
	"github.com/fogleman/gg"
)

const (
	defaultFontSize = 10
	defaultDPI      = 226
)

type Scene struct {
	canvas   *gg.Context
	fontSize float64
}

func NewScene(canvas *gg.Context, opts ...SceneOpt) *Scene {
	s := &Scene{
		canvas:   canvas,
		fontSize: defaultFontSize,
	}
	for _, opt := range opts {
		opt.applyScene(s)
	}
	return s
}
