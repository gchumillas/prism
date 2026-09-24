package ui

import (
	"image/color"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Text struct {
	BasePosition
	scene    *Scene
	text     string
	width    float64
	fontSize float64
}

func (s *Scene) Text(text string, width float64, opts ...TextOpt) *Text {
	t := &Text{
		scene:    s,
		text:     text,
		width:    width,
		fontSize: s.fontSize,
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

func (t *Text) Render() Rectangle {
	face, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    t.fontSize,
		DPI:     defaultDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	m := face.Metrics()
	ascent := float64(m.Ascent) / 64
	lineh := float64(m.Height) / 64

	height := ascent
	lines := splitIntoLines(face, t.text, t.width)
	for _, line := range lines {
		t.scene.canvas.SetFontFace(face)
		t.scene.canvas.SetColor(color.Black)
		// TODO: add ui.Color option
		t.scene.canvas.DrawString(line, float64(t.position.X), float64(t.position.Y+height))
		height += lineh
	}

	return Rectangle{
		Min: t.position,
		Max: t.position.Add(Point{t.width, lineh * float64(len(lines))}),
	}
}

func splitIntoLines(f font.Face, text string, width float64) (lines []string) {
	words := splitIntoWords(text)
	for i, j := 0, 0; i < len(words); {
		var offs float64
		sb := strings.Builder{}
		for k, w := range words[j:] {
			_, advance := font.BoundString(f, w)
			if offs+float64(advance)/64 > width {
				j += k
				break
			}
			offs += float64(advance) / 64
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
