package ui

import (
	_ "embed"
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed assets/RobotoSlab-Regular.ttf
var robotoSlabRegularTTF []byte

var robotoSlabRegular = func() *opentype.Font {
	f, err := opentype.Parse(robotoSlabRegularTTF)
	if err != nil {
		panic(err)
	}
	return f
}()

type Label struct {
	BasePosition
	scene    *Scene
	text     string
	fontSize float64

	// cached results
	cached bool
	face   font.Face
	width  int
	height int
	ascent int
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

func (l *Label) SetPosition(x, y int) {
	l.BasePosition.SetPosition(x, y)
	l.cached = false
}

func (l *Label) Size() (int, int) {
	l.calcSize()
	return l.width, l.height
}

func (l *Label) Render() {
	l.calcSize()
	if l.width <= 0 || l.height <= 0 {
		return
	}
	drawer := &font.Drawer{
		Dst:  l.scene.canvas,
		Src:  image.NewUniform(color.Gray{Y: 0}),
		Face: l.face,
		Dot:  fixed.P(l.X, l.Y+l.ascent),
	}
	drawer.DrawString(l.text)
}

func (l *Label) calcSize() {
	if l.cached {
		return
	}
	var err error
	l.face, err = opentype.NewFace(robotoSlabRegular, &opentype.FaceOptions{
		Size:    l.fontSize,
		DPI:     l.scene.dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	bounds, _ := font.BoundString(l.face, l.text)
	l.width = bounds.Max.X.Ceil()
	metrics := l.face.Metrics()
	l.ascent = metrics.Ascent.Ceil()
	descent := metrics.Descent.Ceil()
	l.height = l.ascent + descent
	l.cached = true
}
