package ui

import (
	"image"
	"reflect"
	"slices"
)

type childState struct {
	pos       image.Point
	regs      []image.Rectangle
	typ       reflect.Type
	footprint []uint64
}

type Stack interface {
	Children() []Widget
	Add(Widget)
	Remove(Widget) int
}

type baseStack struct {
	ctx             Context
	padding         image.Point
	gap             int
	children        []Widget
	prevChildStates []childState
	footprint       func(*baseStack) []uint64
	layout          func(*baseStack, int, int) ([]image.Rectangle, image.Point)
}

type StackOption interface {
	applyStack(*baseStack)
}

func (val WithPadding) applyStack(s *baseStack) { s.padding = image.Point(val) }
func (val WithGap) applyStack(s *baseStack)     { s.gap = int(val) }

func newStack(
	ctx Context,
	footprint func(*baseStack) []uint64,
	layout func(*baseStack, int, int) ([]image.Rectangle, image.Point),
	opts ...StackOption,
) *baseStack {
	s := &baseStack{
		ctx:       ctx,
		footprint: footprint,
		layout:    layout,
	}
	for _, opt := range opts {
		opt.applyStack(s)
	}
	return s
}

func (s *baseStack) Children() []Widget {
	return s.children
}

func (s *baseStack) Add(child Widget) {
	s.children = append(s.children, child)
}

func (s *baseStack) Remove(child Widget) (n int) {
	oldLen := len(s.children)
	s.children = slices.DeleteFunc(s.children, func(c Widget) bool {
		return c == child
	})
	return oldLen - len(s.children)
}

func (s *baseStack) Footprint() []uint64 {
	return s.footprint(s)
}

func (s *baseStack) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	return s.layout(s, x, y)
}
