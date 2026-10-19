package ui

import (
	"errors"
	"fmt"
	"image"
)

var (
	ErrAttached    = errors.New("already attached to another stack")
	ErrNotAttached = errors.New("not found or already detached from the stack")
)

type Button struct {
	wrapper Widget
	widget  Widget
	bounds  image.Rectangle
}

func (b *Button) Contains(x, y int) bool {
	return image.Pt(x, y).In(b.bounds)
}

type clickableWidget struct {
	*Button
}

func NewButton(w Widget) *Button {
	return &Button{widget: w}
}

func (b *Button) AttachTo(parent *Stack) error {
	if b.wrapper != nil {
		return fmt.Errorf("button: %w", ErrAttached)
	}
	b.wrapper = &clickableWidget{Button: b}
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

func (cb *clickableWidget) Render(x, y int) ([]image.Rectangle, image.Point) {
	dirtyRegs, size := cb.widget.Render(x, y)
	pos := image.Pt(x, y)
	cb.bounds = image.Rectangle{Min: pos, Max: pos.Add(size)}
	return dirtyRegs, size
}

func (cb *clickableWidget) Footprint() []uint64 {
	return cb.widget.Footprint()
}
