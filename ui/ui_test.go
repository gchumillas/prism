package ui_test

import (
	"hash/maphash"
	"image"

	"github.com/gchumillas/prism/ui"

	_ "image/png"
)

const (
	appWidth  = 640
	appHeight = 480
)

var seed = maphash.MakeSeed()

type Env struct {
	canvas *image.RGBA
	ctx    ui.Context
}

func NewEnv() *Env {
	canvas := image.NewRGBA(image.Rect(0, 0, appWidth, appHeight))
	ctx := ui.NewContextFromRGBA(canvas)
	return &Env{canvas: canvas, ctx: ctx}
}
