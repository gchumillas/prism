package ui

import (
	"fmt"
	"image"
)

type buttonWrapper struct {
	*Button
}

func (cb *buttonWrapper) Render(x, y int) []image.Rectangle {
	dirtyRegs := cb.widget.Render(x, y)
	pos := image.Pt(x, y)
	cb.bounds = image.Rectangle{Min: pos, Max: pos.Add(cb.Size())}
	return dirtyRegs
}

func (cb *buttonWrapper) Size() image.Point {
	return cb.widget.Size()
}

func (cb *buttonWrapper) Footprint() []uint64 {
	return cb.widget.Footprint()
}

type Button struct {
	wrapper Widget
	widget  Widget
	bounds  image.Rectangle
}

func (b *Button) Contains(x, y int) bool {
	return image.Pt(x, y).In(b.bounds)
}

func NewButton(w Widget) *Button {
	return &Button{widget: w}
}

func (b *Button) AttachTo(parent *Stack) error {
	if b.wrapper != nil {
		return fmt.Errorf("button: %w", ErrAttached)
	}
	b.wrapper = &buttonWrapper{Button: b}
	parent.Add(b.wrapper)
	return nil
}

func (b *Button) DetachFrom(parent *Stack) error {
	if parent.Remove(b.wrapper) > 0 {
		b.wrapper = nil
		return nil
	}
	return ErrNotAttached
}
