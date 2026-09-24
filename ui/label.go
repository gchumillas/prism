package ui

import (
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Label struct {
	BasePosition
	scene    *Scene
	text     string
	fontSize float64
}

func (s *Scene) Label(text string, opts ...LabelOpt) *Label {
	l := &Label{
		scene:    s,
		fontSize: s.fontSize,
		text:     text,
	}
	for _, opt := range opts {
		opt.applyLabel(l)
	}
	return l
}

func (l *Label) Height() float64 {
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
	ascent := float64(metrics.Ascent) / 64
	descent := float64(metrics.Descent) / 64
	return ascent + descent
}

func (l *Label) Width() float64 {
	f, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	bounds, _ := font.BoundString(f, l.text)
	return float64(bounds.Max.X) / 64
}

func (l *Label) Size() Point {
	f, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	bounds, _ := font.BoundString(f, l.text)
	width := float64(bounds.Max.X) / 64

	metrics := f.Metrics()
	ascent := float64(metrics.Ascent) / 64
	descent := float64(metrics.Descent) / 64
	height := ascent + descent

	return Point{X: width, Y: height}
}

func (l *Label) Render() (r Rectangle) {
	f, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	bounds, _ := font.BoundString(f, l.text)
	w := float64(bounds.Max.X) / 64

	metrics := f.Metrics()
	ascent := float64(metrics.Ascent) / 64
	descent := float64(metrics.Descent) / 64
	h := ascent + descent

	l.scene.canvas.SetFontFace(f)
	l.scene.canvas.SetColor(color.Black)
	// TODO: add ui.Color option
	l.scene.canvas.DrawString(l.text, float64(l.position.X), l.position.Y+float64(ascent))

	size := Point{w, h}
	return Rectangle{
		Min: l.position,
		Max: l.position.Add(size),
	}
}
