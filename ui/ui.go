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
	Render() image.Rectangle
}
