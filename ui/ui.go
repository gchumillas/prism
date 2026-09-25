package ui

import (
	"image"

	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
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
	// Render draws the element and returns the dirty region that the
	// refresh system must repaint. This may be larger than the element's
	// own bounds (e.g. to cover borders or effects drawn outside them).
	Render() image.Rectangle
}
