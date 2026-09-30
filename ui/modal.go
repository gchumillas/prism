package ui

import "image"

type Modal struct {
	BaseGroup
	BasePosition
	BaseSize
}

func (c Context) Modal() *Modal {
	m := &Modal{}
	r := c.canvas.Image().Bounds()
	t := image.Rectangle{
		Min: image.Point{X: r.Min.X, Y: r.Min.Y},
		Max: image.Point{X: r.Max.X, Y: r.Max.Y},
	}
	p := t.Max.Sub(t.Min)
	m.position = image.Point{(p.X - 300) / 2, (p.Y - 300) / 2}
	return m
}

func (m *Modal) Render() (rect []image.Rectangle) {
	for _, c := range m.children {
		c.Render()
	}
	return
}
