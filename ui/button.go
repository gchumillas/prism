package ui

import (
	"image"
)

type Button struct {
	wrapper Widget
	bounds  image.Rectangle
}

func (b *Button) Contains(x, y int) bool {
	return image.Pt(x, y).In(b.bounds)
}

type clickableWidget struct {
	*Button
	widget Widget
}

func AttachButton(owner Stack, w Widget) *Button {
	btn := &Button{}
	btn.wrapper = &clickableWidget{Button: btn, widget: w}
	owner.Add(btn.wrapper)
	return btn
}

func (b *Button) DetachFrom(owner Stack) bool {
	return owner.Remove(b.wrapper) > 0
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
