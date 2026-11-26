package ui

import (
	"hash/maphash"
	"image"
	"image/color"
	"image/draw"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Label struct {
	canvas   draw.Image
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

func NewLabel(canvas draw.Image, text string, opts ...LabelOption) *Label {
	l := &Label{
		canvas:   canvas,
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

func (l *Label) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	d := &font.Drawer{Dst: l.canvas, Src: image.NewUniform(l.color), Face: l.face, Dot: fixed.P(pos.X, pos.Y+l.ascent)}
	d.DrawString(l.text)

	return []image.Rectangle{{Min: pos, Max: pos.Add(l.Size())}}
}
