package ui

import (
	"hash/maphash"
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Label struct {
	canvas   *gg.Context
	text     string
	fontSize float64
	color    color.NRGBA

	face    font.Face
	descent int
	ascent  int
}

type LabelOption interface {
	ApplyLabel(*Label)
}

func (val WithFontSize) ApplyLabel(l *Label) { l.fontSize = float64(val) }
func (val WithAlpha) ApplyLabel(l *Label)    { l.color.A = uint8(val) }
func (val WithColor) ApplyLabel(l *Label) {
	v := uint32(val)
	l.color.R = uint8(v >> 16)
	l.color.G = uint8(v >> 8)
	l.color.B = uint8(v)
}

func NewLabel(c *Canvas, text string, opts ...LabelOption) *Label {
	l := &Label{
		canvas:   c.context,
		fontSize: defaultFontSize,
		color:    defaultColor,
		text:     text,
	}
	for _, opt := range opts {
		opt.ApplyLabel(l)
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

func (l *Label) Footprint() []uint64 {
	r, g, b, a := l.color.R, l.color.G, l.color.B, l.color.A
	c := uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a)
	return []uint64{math.Float64bits(l.fontSize), uint64(c), maphash.String(seed, l.text)}
}

func (l *Label) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	l.canvas.SetFontFace(l.face)
	l.canvas.SetColor(l.color)
	l.canvas.DrawString(l.text, float64(pos.X), float64(pos.Y+l.ascent))

	bounds, _ := font.BoundString(l.face, l.text)
	size = image.Pt(bounds.Max.X.Ceil(), l.descent+l.ascent)

	return []image.Rectangle{{Min: pos, Max: pos.Add(size)}}, size
}
