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
	m.SetSize(300, 300)
	r := s.canvas.Bounds()
	p := r.Max.Sub(r.Min)
	m.x, m.y = (p.X-m.width)/2, (p.Y-m.height)/2
	return m
}

func (m *Modal) Render() (rect image.Rectangle) {
	r := image.Rect(m.x, m.y, m.x+m.width, m.y+m.height)
	draw.Draw(m.scene.canvas, r, image.NewUniform(color.White), image.Pt(0, 0), draw.Src)
	drawRect(m.scene.canvas, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color.Black)

	for _, c := range m.children {
		c.Render()
	}
	return
}
