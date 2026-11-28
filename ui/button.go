package ui

import (
	"fmt"
	"image"
)

type buttonWrapper[T Widget] struct {
	*Button[T]
}

func (cb *buttonWrapper[T]) Render(x, y int) []image.Rectangle {
	dirtyRegs := cb.widget.Render(x, y)
	pos := image.Pt(x, y)
	cb.bounds = image.Rectangle{Min: pos, Max: pos.Add(cb.Size())}
	return dirtyRegs
}

func (cb *buttonWrapper[T]) Size() image.Point {
	return cb.widget.Size()
}

func (cb *buttonWrapper[T]) Footprint() []uint64 {
	return cb.widget.Footprint()
}

type Button[T Widget] struct {
	wrapper Widget
	widget  T
	bounds  image.Rectangle
}

func NewButton[T Widget](w T) *Button[T] {
	return &Button[T]{widget: w}
}

func (b *Button[T]) Contains(x, y int) bool {
	return image.Pt(x, y).In(b.bounds)
}

func (b *Button[T]) Widget() T {
	return b.widget
}

func (b *Button[T]) AttachTo(parent *Stack) error {
	if b.wrapper != nil {
		return fmt.Errorf("button: %w", ErrAttached)
	}
	b.wrapper = &buttonWrapper[T]{Button: b}
	parent.Add(b.wrapper)
	return nil
}

func (b *Button[T]) DetachFrom(parent *Stack) error {
	if parent.Remove(b.wrapper) > 0 {
		parent.outdated = true
		b.wrapper = nil
		return nil
	}
	return ErrNotAttached
}
