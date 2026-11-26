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
	color    color.NRGBA
}

type TextOption interface {
	ApplyText(*Text)
}

func (val WithFontSize) ApplyText(l *Text) { l.fontSize = float64(val) }
func (val WithAlpha) ApplyText(l *Text)    { l.color.A = uint8(val) }
func (val WithColor) ApplyText(l *Text) {
	v := uint32(val)
	l.color.R = uint8(v >> 16)
	l.color.G = uint8(v >> 8)
	l.color.B = uint8(v)
}

func NewText(canvas draw.Image, text string, width int, opts ...TextOption) *Text {
	t := &Text{
		canvas:   canvas,
		text:     text,
		width:    width,
		fontSize: defaultFontSize,
		color:    defaultColor,
	}
	for _, opt := range opts {
		opt.ApplyText(t)
	}
	return t
}

func (t *Text) Footprint() []uint64 {
	r, g, b, a := t.color.R, t.color.G, t.color.B, t.color.A
	c := uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a)
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

	d := &font.Drawer{Dst: t.canvas, Src: image.NewUniform(t.color), Face: face}
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
