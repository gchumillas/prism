package ui

import (
	"image"
	"image/color"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Text struct {
	BasePosition
	canvas   *gg.Context
	text     string
	width    int
	fontSize float64
}

type TextOpt interface {
	applyText(*Text)
}

// avail options
func (val xOpt) applyText(x *Text)        { x.SetX(int(val)) }
func (val yOpt) applyText(x *Text)        { x.SetY(int(val)) }
func (val fontSizeOpt) applyText(x *Text) { x.fontSize = float64(val) }

func (c Context) Text(text string, width int, opts ...TextOpt) *Text {
	t := &Text{
		canvas:   c.canvas,
		text:     text,
		width:    width,
		fontSize: c.fontSize,
	}
	for _, opt := range opts {
		opt.applyText(t)
	}
	return t
}

func (t *Text) Size() (width, height int) {
	return
}

func (bs *Text) SetSize(width, height int) {
}

func (t *Text) Render() image.Rectangle {
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

	height := ascent
	lines := splitIntoLines(face, t.text, t.width)
	for _, line := range lines {
		t.canvas.SetFontFace(face)
		t.canvas.SetColor(color.Black)
		// TODO: add ui.Color option
		t.canvas.DrawString(line, float64(t.position.X), float64(t.position.Y+height))
		height += lineh
	}

	return image.Rectangle{
		Min: t.position,
		Max: t.position.Add(image.Pt(t.width, lineh*len(lines))),
	}
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
