package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Label struct {
	BasePosition
	canvas   *gg.Context
	text     string
	fontSize float64

	face    font.Face
	descent int
	ascent  int
}

type LabelOpt interface {
	applyLabel(*Label)
}

// avail options
func (val fontSizeOpt) applyLabel(x *Label) { x.fontSize = float64(val) }

func (c *Context) Label(text string, opts ...LabelOpt) *Label {
	l := &Label{
		canvas:   c.canvas,
		fontSize: c.fontSize,
		text:     text,
	}
	for _, opt := range opts {
		opt.applyLabel(l)
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

func (l *Label) Height() int {
	return l.descent + l.ascent
}

func (l *Label) Width() int {
	bounds, _ := font.BoundString(l.face, l.text)
	return bounds.Max.X.Ceil()
}

func (l *Label) Render() (r image.Rectangle) {
	l.canvas.SetFontFace(l.face)
	l.canvas.SetColor(color.Black)
	l.canvas.DrawString(l.text, float64(l.position.X), float64(l.position.Y+l.ascent))

	size := image.Point{l.Width(), l.Height()}
	return image.Rectangle{
		Min: l.position,
		Max: l.position.Add(size),
	}
}
