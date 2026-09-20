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
	scene    *Scene
	text     string
	fontSize float64
	x, y     int

	// cached results
	cached bool
	face   font.Face
	width  int
	height int
	ascent int
}

func (s *Scene) Label(text string) *Label {
	return &Label{
		scene:    s,
		fontSize: s.fontSize,
		text:     text,
	}
}

func (t *Label) FontSize() float64 {
	return t.fontSize
}

func (t *Label) SetFontSize(val float64) {
	t.fontSize = val
	t.cached = false
}

func (t *Label) Position() (int, int) {
	return t.x, t.y
}

func (t *Label) SetPosition(x, y int) {
	t.x, t.y = x, y
	t.cached = false
}

func (t *Label) Size() (int, int) {
	t.calcSize()
	return t.width, t.height
}

func (t *Label) Render() {
	t.calcSize()
	if t.width <= 0 || t.height <= 0 {
		return
	}
	drawer := &font.Drawer{
		Dst:  t.scene.canvas,
		Src:  image.NewUniform(color.Gray{Y: 0}),
		Face: t.face,
		Dot:  fixed.P(t.x, t.y+t.ascent),
	}
	drawer.DrawString(t.text)
}

func (t *Label) calcSize() {
	if t.cached {
		return
	}
	var err error
	t.face, err = opentype.NewFace(robotoSlabRegular, &opentype.FaceOptions{
		Size:    t.fontSize,
		DPI:     t.scene.dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	bounds, _ := font.BoundString(t.face, t.text)
	t.width = bounds.Max.X.Ceil()
	metrics := t.face.Metrics()
	t.ascent = metrics.Ascent.Ceil()
	descent := metrics.Descent.Ceil()
	t.height = t.ascent + descent
	t.cached = true
}
