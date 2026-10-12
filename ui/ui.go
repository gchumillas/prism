package ui

import (
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

type Widget interface {
	Render(x, y int) (dirtyRegs []image.Rectangle)
	Size() image.Point
	Footprint() []any
}

type Context struct {
	canvas   *gg.Context
	FontSize float64
}

func NewContext(width, height int) Context {
	canvas := gg.NewContext(width, height)
	return Context{
		canvas:   canvas,
		FontSize: defaultFontSize,
	}
}

func NewContextFromRGBA(img *image.RGBA) Context {
	canvas := gg.NewContextForRGBA(img)
	return Context{
		canvas:   canvas,
		FontSize: defaultFontSize,
	}
}

func (c Context) Canvas() image.Image {
	return c.canvas.Image()
}
