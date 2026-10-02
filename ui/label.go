package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Label struct {
	Position image.Point
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

func (l *Label) SetPosition(val image.Point) {
	l.Position = val
}

func (l *Label) Height() int {
	return l.descent + l.ascent
}

func (l *Label) Width() int {
	bounds, _ := font.BoundString(l.face, l.text)
	return bounds.Max.X.Ceil()
}

func (l *Label) Render() []image.Rectangle {
	l.canvas.SetFontFace(l.face)
	l.canvas.SetColor(color.Black)
	l.canvas.DrawString(l.text, float64(l.Position.X), float64(l.Position.Y+l.ascent))

	size := image.Point{l.Width(), l.Height()}
	return []image.Rectangle{{
		Min: l.Position,
		Max: l.Position.Add(size),
	}}
}
