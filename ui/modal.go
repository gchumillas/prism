package ui

type Modal struct {
	BaseGroup
	BasePosition
	BaseSize
	scene *Scene
}

func (s *Scene) Modal() *Modal {
	m := &Modal{scene: s}
	r := s.canvas.Image().Bounds()
	t := Rectangle{
		Min: Point{X: float64(r.Min.X), Y: float64(r.Min.Y)},
		Max: Point{X: float64(r.Max.X), Y: float64(r.Max.Y)},
	}
	p := t.Max.Sub(t.Min)
	m.position = Point{(p.X - 300) / 2, (p.Y - 300) / 2}
	return m
}

func (m *Modal) Render() (rect Rectangle) {
	// r := image.Rect(m.position.X, m.position.Y, m.position.X+m.size.X, m.position.Y+m.size.Y)
	// img := m.scene.canvas.Image()
	// draw.Draw(img, r, image.NewUniform(color.White), image.Pt(0, 0), draw.Src)
	// drawRect(img, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color.Black)

	for _, c := range m.children {
		c.Render()
	}
	return
}
