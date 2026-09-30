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

type Renderer interface {
	Render() (dirty []image.Rectangle)
}

type Context struct {
	canvas   *gg.Context
	fontSize float64
}

type ContextOpt interface {
	applyContext(*Context)
}

// avail options
func (val fontSizeOpt) applyContext(x *Context) { x.fontSize = float64(val) }

func NewContext(width, height int, opts ...ContextOpt) Context {
	canvas := gg.NewContext(width, height)
	c := &Context{
		canvas:   canvas,
		fontSize: defaultFontSize,
	}
	for _, opt := range opts {
		opt.applyContext(c)
	}
	return *c
}

func NewContextFromRGBA(img *image.RGBA, opts ...ContextOpt) Context {
	canvas := gg.NewContextForRGBA(img)
	c := &Context{
		canvas:   canvas,
		fontSize: defaultFontSize,
	}
	for _, opt := range opts {
		opt.applyContext(c)
	}
	return *c
}

func (c Context) Image() image.Image {
	return c.canvas.Image()
}
