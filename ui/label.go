package ui

import (
	"hash/maphash"
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

type labelOption interface {
	applyLabel(*Label)
}

func (val WithFontSize) applyLabel(l *Label) { l.fontSize = float64(val) }

// TODO: add optional color param
func NewLabel(c Context, text string, opts ...labelOption) *Label {
	l := &Label{
		canvas:   c.canvas,
		fontSize: defaultFontSize,
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

func (l *Label) size() image.Point {
	bounds, _ := font.BoundString(l.face, l.text)
	return image.Pt(bounds.Max.X.Ceil(), l.descent+l.ascent)
}

func (l *Label) Footprint() []uint64 {
	// TODO: maybe we don'tneed maphash (affects other places)
	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, l.text)
	return []uint64{h.Sum64()}
}

func (l *Label) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	l.canvas.SetFontFace(l.face)
	l.canvas.SetColor(color.Black)
	l.canvas.DrawString(l.text, float64(pos.X), float64(pos.Y+l.ascent))

	bounds, _ := font.BoundString(l.face, l.text)
	size = image.Pt(bounds.Max.X.Ceil(), l.descent+l.ascent)

	return []image.Rectangle{{Min: pos, Max: pos.Add(size)}}, size
}
