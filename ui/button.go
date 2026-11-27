package ui

import (
	"fmt"
	"image"
	"slices"
)

type Button[T Widget] struct {
	widget   T
	bounds   image.Rectangle
	attached bool
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
	if b.attached {
		return fmt.Errorf("button: %w", ErrAttached)
	}
	parent.children = append(parent.children, b)
	b.attached = true
	return nil
}

func (b *Button[T]) DetachFrom(parent *Stack) error {
	if i := slices.IndexFunc(parent.children, func(x Drawable) bool {
		return x == b
	}); i > -1 {
		parent.children = slices.Delete(parent.children, i, i+1)
		b.attached = false
		parent.outdated = true
		return nil
	}
	return ErrNotAttached
}

func (b *Button[T]) Render(x, y int) []image.Rectangle {
	dirtyRegs := b.widget.Render(x, y)
	pos := image.Pt(x, y)
	b.bounds = image.Rectangle{Min: pos, Max: pos.Add(b.widget.Size())}
	return dirtyRegs
}

func (b *Button[T]) Footprint() []uint64 {
	return b.widget.Footprint()
}

func (b *Button[T]) Size() image.Point {
	return b.widget.Size()
}
