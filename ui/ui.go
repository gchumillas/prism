package ui

import (
	"errors"
	"hash/maphash"
	"image"
	"image/color"

	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

var (
	defaultFontSize = 10.
	defaultColor    = color.NRGBA{A: 0xff}
	defaultDPI      = 226.
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

var (
	ErrAttached    = errors.New("already attached to another stack")
	ErrNotAttached = errors.New("not found or already detached from the stack")
	ErrCircularRef = errors.New("circular reference")
	ErrNeedsRender = errors.New("call Render() first")
)

type Widget interface {
	Render(x, y int) (dirtyRegs []image.Rectangle)
	Size() image.Point
	Footprint() []uint64
}
