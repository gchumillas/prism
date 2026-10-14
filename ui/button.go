package ui

import "image"

type Button struct {
	bounds image.Rectangle
}

func (b *Button) Contains(x, y int) bool {
	return image.Pt(x, y).In(b.bounds)
}

type clickableWidget struct {
	*Button
	widget Widget
}

func newClickableWidget(w Widget) *clickableWidget {
	return &clickableWidget{Button: &Button{}, widget: w}
}

func (cb *clickableWidget) Render(x, y int) ([]image.Rectangle, image.Point) {
	dirtyRegs, size := cb.widget.Render(x, y)
	pos := image.Pt(x, y)
	cb.bounds = image.Rectangle{Min: pos, Max: pos.Add(size)}
	return dirtyRegs, size
}

func (cb *clickableWidget) Footprint() []uint64 {
	return cb.widget.Footprint()
}
