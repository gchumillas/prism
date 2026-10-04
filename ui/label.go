package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Label struct {
	canvas   *gg.Context
	text     string
	fontSize float64

	face    font.Face
	descent int
	ascent  int
}

func NewLabel(c Context, text string) *Label {
	l := &Label{
		canvas:   c.canvas,
		fontSize: c.FontSize,
		text:     text,
	}

	var err error
	l.face, err = opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	metrics := l.face.Metrics()
	l.descent = metrics.Descent.Ceil()
	l.ascent = metrics.Ascent.Ceil()
	return l
}

func (l *Label) Size() image.Point {
	bounds, _ := font.BoundString(l.face, l.text)
	return image.Pt(bounds.Max.X.Ceil(), l.descent+l.ascent)
}

func (l *Label) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	l.canvas.SetFontFace(l.face)
	l.canvas.SetColor(color.Black)
	l.canvas.DrawString(l.text, float64(pos.X), float64(pos.Y+l.ascent))

	size := l.Size()
	return []image.Rectangle{{
		Min: pos,
		Max: pos.Add(size),
	}}
}
