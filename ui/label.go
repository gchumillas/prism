package ui

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
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

func (l *Label) Render() (r image.Rectangle) {
	f, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     l.scene.dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	m := f.Metrics()
	ascent := m.Ascent.Ceil()

	drawer := &font.Drawer{
		Dst: l.scene.canvas,
		// TODO: add ui.Color option
		Src:  image.NewUniform(color.Black),
		Face: f,
		Dot:  fixed.P(l.position.X, l.position.Y+ascent),
	}
	drawer.DrawString(l.text)

	bounds, _ := font.BoundString(f, l.text)
	w := bounds.Max.X.Ceil()

	metrics := f.Metrics()
	ascent = metrics.Ascent.Ceil()
	descent := metrics.Descent.Ceil()
	h := ascent + descent

	size := image.Pt(w, h)
	return image.Rectangle{
		Min: l.position,
		Max: l.position.Add(size),
	}
}
