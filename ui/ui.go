package ui

import (
	"hash/maphash"
	"image"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

const (
	defaultFontSize = 10
	defaultDPI      = 226
)

// add ui.Font() option
var uiFont = func() *opentype.Font {
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	return f
}()

var seed = maphash.MakeSeed()

type Widget interface {
	Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point)
	Footprint() []uint64
}

type Canvas struct {
	context *gg.Context
}

func NewCanvas(width, height int) *Canvas {
	dc := gg.NewContext(width, height)
	return &Canvas{context: dc}
}

func NewCanvasFromRGBA(img *image.RGBA) *Canvas {
	dc := gg.NewContextForRGBA(img)
	return &Canvas{context: dc}
}

func (c *Canvas) Image() image.Image {
	return c.context.Image()
}
