package ui

import (
	"fmt"
	"image"
	"slices"
)

type Button struct {
	widget   Widget
	bounds   image.Rectangle
	attached bool
}

func NewButton(w Widget) *Button {
	return &Button{widget: w}
}

func (b *Button) Contains(x, y int) bool {
	return image.Pt(x, y).In(b.bounds)
}

func (b *Button) AttachTo(parent *Stack) error {
	if b.attached {
		return fmt.Errorf("button: %w", ErrAttached)
	}
	parent.children = append(parent.children, b)
	b.attached = true
	return nil
}

func (b *Button) DetachFrom(parent *Stack) error {
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

func (b *Button) Render(x, y int) []image.Rectangle {
	dirtyRegs := b.widget.Render(x, y)
	pos := image.Pt(x, y)
	b.bounds = image.Rectangle{Min: pos, Max: pos.Add(b.widget.Size())}
	return dirtyRegs
}

func (b *Button) Size() image.Point {
	return b.widget.Size()
}
