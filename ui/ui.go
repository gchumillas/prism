package ui

import (
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

var uiFont = func() *opentype.Font {
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	return f
}()

type Renderer interface {
	Render()
}
