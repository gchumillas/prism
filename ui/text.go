package ui

import (
	"image"
	"image/color"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Text struct {
	BasePosition
	scene    *Scene
	text     string
	width    int
	fontSize float64
}

func (s *Scene) Text(text string, width int, opts ...TextOpt) *Text {
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

func (t *Text) Render() (r image.Rectangle) {
	face, err := opentype.NewFace(uiFont, &opentype.FaceOptions{
		Size:    t.fontSize,
		DPI:     t.scene.dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	m := face.Metrics()
	ascent := m.Ascent.Ceil()
	height := m.Height.Ceil()

	lines := splitIntoLines(face, t.text, t.width)
	for i, line := range lines {
		d := &font.Drawer{
			Dst:  t.scene.canvas,
			Src:  image.NewUniform(color.Gray{Y: 0}),
			Face: face,
			Dot:  fixed.P(t.position.X, t.position.Y+ascent+height*i),
		}
		d.DrawString(line)
	}
	return
}

func splitIntoLines(f font.Face, text string, width int) (lines []string) {
	words := splitIntoWords(text)
	maxWidth := fixed.I(width)
	var offset fixed.Int26_6
	var accWords []string
	for _, w := range words {
		_, advance := font.BoundString(f, w)
		if offset+advance > maxWidth {
			if len(accWords) > 0 {
				lines = append(lines, strings.Join(accWords, " "))
			}
			accWords = nil
			offset = 0
		}
		accWords = append(accWords, w)
		offset += advance
	}
	if len(accWords) > 0 {
		lines = append(lines, strings.Join(accWords, " "))
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
