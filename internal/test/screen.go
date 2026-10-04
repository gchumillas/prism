package test

import (
	"image"
	"image/draw"
)

type Screen struct {
	Canvas  *image.RGBA // internal canvas
	Display *image.RGBA // physical display
}

func NewScreen(w, h int) *Screen {
	r := image.Rect(0, 0, w, h)
	return &Screen{
		Canvas:  image.NewRGBA(r),
		Display: image.NewRGBA(r),
	}
}

func (s *Screen) Refresh() {
	r := s.Canvas.Bounds()
	s.RefreshRecs([]image.Rectangle{r})
}

func (s *Screen) RefreshRecs(recs []image.Rectangle) {
	for _, r := range recs {
		draw.Draw(s.Display, r, s.Canvas, r.Min, draw.Src)
	}
}
