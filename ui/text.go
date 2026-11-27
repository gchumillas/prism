package ui

import (
	"hash/maphash"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Text struct {
	canvas   draw.Image
	text     string
	width    int
	fontSize float64
	Color    color.Color
}

func NewText(canvas draw.Image, text string, fontSize float64, width int) *Text {
	t := &Text{
		canvas:   canvas,
		text:     text,
		width:    width,
		fontSize: fontSize,
		Color:    defaultColor,
	}
	return t
}

func (t *Text) Footprint() []uint64 {
	r, g, b, a := t.Color.RGBA()
	c := uint64(r)<<48 | uint64(g)<<32 | uint64(b)<<16 | uint64(a)
	return []uint64{math.Float64bits(t.fontSize), uint64(c), uint64(t.width), maphash.String(seed, t.text)}
}

func (t *Text) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	face, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    t.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	m := face.Metrics()
	ascent := m.Ascent.Ceil()
	lineh := m.Height.Ceil()

	d := &font.Drawer{Dst: t.canvas, Src: image.NewUniform(t.Color), Face: face}
	h := ascent
	lines := splitIntoLines(face, t.text, t.width)
	for _, line := range lines {
		d.Dot = fixed.P(pos.X, pos.Y+h)
		d.DrawString(line)
		h += lineh
	}

	return []image.Rectangle{{
		Min: pos,
		Max: pos.Add(image.Pt(t.width, lineh*len(lines))),
	}}
}

func (t *Text) Size() image.Point {
	face, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    t.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	m := face.Metrics()
	lineh := m.Height.Ceil()
	lines := splitIntoLines(face, t.text, t.width)
	return image.Pt(t.width, lineh*len(lines))
}

func splitIntoLines(f font.Face, text string, width int) (lines []string) {
	words := splitIntoWords(text)
	for i, j := 0, 0; i < len(words); {
		var offs fixed.Int26_6
		sb := strings.Builder{}
		for k, w := range words[j:] {
			_, advance := font.BoundString(f, w)
			if offs+advance > fixed.I(width) {
				j += k
				break
			}
			offs += advance
			sb.WriteString(w)
			i++
		}
		line := sb.String()
		if j > 0 {
			line = strings.TrimLeft(line, " ")
		}
		lines = append(lines, line)
	}
	return
}

func splitIntoWords(text string) (words []string) {
	i, count := 0, utf8.RuneCountInString(text)
	for i < count {
		var s strings.Builder
		// spaces
		for {
			r, size := utf8.DecodeRuneInString(text[i:])
			if size == 0 || !unicode.IsSpace(r) {
				break
			}
			i += size
			s.WriteRune(r)
		}
		// letters
		for {
			r, size := utf8.DecodeRuneInString(text[i:])
			if size == 0 || unicode.IsSpace(r) {
				break
			}
			i += size
			s.WriteRune(r)
		}
		words = append(words, s.String())
	}
	return
}
