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
	r := s.canvas.Bounds()
	p := r.Max.Sub(r.Min)
	m.position = image.Pt((p.X-300)/2, (p.Y-300)/2)
	return m
}

func (m *Modal) Render() (rect image.Rectangle) {
	r := image.Rect(m.position.X, m.position.Y, m.position.X+m.width, m.position.Y+m.height)
	draw.Draw(m.scene.canvas, r, image.NewUniform(color.White), image.Pt(0, 0), draw.Src)
	drawRect(m.scene.canvas, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color.Black)

	for _, c := range m.children {
		c.Render()
	}
	return
}
