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
}

type LabelOpt interface {
	applyLabel(*Label)
}

// avail options
func (val fontSizeOpt) applyLabel(x *Label) { x.fontSize = float64(val) }

func (s *Scene) Label(text string, opts ...LabelOpt) *Label {
	l := &Label{
		canvas:   s.canvas,
		fontSize: s.fontSize,
		text:     text,
	}
	for _, opt := range opts {
		opt.applyLabel(l)
	}
	return l
}

func (l *Label) Height() int {
	// TODO: implement some caching system to avoid unnecessary recalculation
	f, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	metrics := f.Metrics()
	ascent := metrics.Ascent.Ceil()
	descent := metrics.Descent.Ceil()
	return ascent + descent
}

func (l *Label) Width() int {
	f, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	bounds, _ := font.BoundString(f, l.text)
	return bounds.Max.X.Ceil()
}

func (l *Label) Size() image.Point {
	f, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	bounds, _ := font.BoundString(f, l.text)
	width := bounds.Max.X.Ceil()

	metrics := f.Metrics()
	ascent := metrics.Ascent.Ceil()
	descent := metrics.Descent.Ceil()
	height := ascent + descent

	return image.Point{X: width, Y: height}
}

func (l *Label) Render() (r image.Rectangle) {
	f, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	bounds, _ := font.BoundString(f, l.text)
	w := bounds.Max.X.Ceil()

	metrics := f.Metrics()
	ascent := metrics.Ascent.Ceil()
	descent := metrics.Descent.Ceil()
	h := ascent + descent

	l.canvas.SetFontFace(f)
	l.canvas.SetColor(color.Black)
	// TODO: add ui.Color option
	l.canvas.DrawString(l.text, float64(l.position.X), float64(l.position.Y+ascent))

	size := image.Point{w, h}
	return image.Rectangle{
		Min: l.position,
		Max: l.position.Add(size),
	}
}
