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

type Stack struct {
	ctx             Context
	padding         image.Point
	gap             int
	children        []Widget
	prevChildStates []childState
	footprint       func(*Stack) []uint64
	layout          func(*Stack, int, int) ([]image.Rectangle, image.Point)
}

type StackOption interface {
	applyStack(*Stack)
}

func (val WithPadding) applyStack(s *Stack) { s.padding = image.Point(val) }
func (val WithGap) applyStack(s *Stack)     { s.gap = int(val) }

func newStack(
	ctx Context,
	footprint func(*Stack) []uint64,
	layout func(*Stack, int, int) ([]image.Rectangle, image.Point),
	opts ...StackOption,
) *Stack {
	s := &Stack{
		ctx:       ctx,
		footprint: footprint,
		layout:    layout,
	}
	for _, opt := range opts {
		opt.applyStack(s)
	}
	return s
}

// TODO: remove Stack.Children
func (s *Stack) Children() []Widget {
	return s.children
}

func (s *Stack) Add(child Widget) {
	s.children = append(s.children, child)
}

func (s *Stack) Remove(child Widget) (n int) {
	oldLen := len(s.children)
	s.children = slices.DeleteFunc(s.children, func(c Widget) bool {
		return c == child
	})
	return oldLen - len(s.children)
}

func (s *Stack) AddButton(child Widget) *Button {
	btn := newClickableWidget(child)
	s.children = append(s.children, btn)
	return btn.Button
}

func (s *Stack) RemoveButton(btn *Button) bool {
	i := slices.IndexFunc(s.children, func(x Widget) bool {
		if v, ok := x.(*clickableWidget); ok {
			if v.Button == btn {
				return true
			}
		}
		return false
	})
	if i > -1 {
		s.children = slices.Delete(s.children, i, i+1)
		return true
	}
	return false
}

func (s *Stack) Footprint() []uint64 {
	return s.footprint(s)
}

func (s *Stack) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	return s.layout(s, x, y)
}
