package ui

import (
	"image"
	"image/color"
	"image/draw"
)

type Modal struct {
	BaseGroup
	BasePosition
	BaseSize
	scene *Scene
}

func (s *Scene) Modal() *Modal {
	m := &Modal{scene: s}
	m.Width, m.Height = 300, 300
	r := s.canvas.Bounds()
	p := r.Max.Sub(r.Min)
	m.SetPosition((p.X-m.Width)/2, (p.Y-m.Width)/2)
	return m
}

func (m *Modal) Render() {
	x, y := m.Position()
	r := image.Rect(x, y, x+m.Width, y+m.Height)
	draw.Draw(m.scene.canvas, r, image.NewUniform(color.White), image.Pt(0, 0), draw.Src)
	drawRect(m.scene.canvas, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color.Black)

	for _, c := range m.children {
		c.Render()
	}
}
