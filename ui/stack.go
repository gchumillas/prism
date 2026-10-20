package ui

import (
	"errors"
	"fmt"
	"image"
	"reflect"
	"slices"
)

var ErrCircularRef = errors.New("circular reference")

type childState struct {
	pos       image.Point
	regs      []image.Rectangle
	typ       reflect.Type
	footprint []uint64
}

type stackWrapper struct {
	*Stack
}

func (sw *stackWrapper) Footprint() []uint64 {
	return sw.footprint()
}

type Stack struct {
	wrapper         *stackWrapper
	canvas          *Canvas
	paddingX        int
	paddingY        int
	gap             int
	children        []Widget
	prevChildStates []childState
	stamp           func(*Stack) []uint64
	layout          func(*Stack, int, int) ([]image.Rectangle, image.Point)
}

type StackOption interface {
	ApplyStack(*Stack)
}

func (val WithPaddingX) ApplyStack(s *Stack) { s.paddingX = int(val) }
func (val WithPaddingY) ApplyStack(s *Stack) { s.paddingY = int(val) }
func (val WithGap) ApplyStack(s *Stack)      { s.gap = int(val) }

func newStack(
	canvas *Canvas,
	footprint func(*Stack) []uint64,
	layout func(*Stack, int, int) ([]image.Rectangle, image.Point),
	opts ...StackOption,
) *Stack {
	s := &Stack{
		canvas: canvas,
		stamp:  footprint,
		layout: layout,
	}
	for _, opt := range opts {
		opt.ApplyStack(s)
	}
	return s
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

func (s *Stack) AttachTo(parent *Stack) error {
	if s.wrapper != nil {
		return fmt.Errorf("stack: %w", ErrAttached)
	}
	if parent == s || s.hasStack(parent) {
		return fmt.Errorf("stack: %w", ErrCircularRef)
	}
	s.wrapper = &stackWrapper{Stack: s}
	parent.Add(s.wrapper)
	return nil
}

func (s *Stack) DetachFrom(parent *Stack) error {
	if parent.Remove(s.wrapper) > 0 {
		s.wrapper = nil
		return nil
	}
	return ErrNotAttached
}

func (s *Stack) footprint() []uint64 {
	return s.stamp(s)
}

func (s *Stack) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	return s.layout(s, x, y)
}

func (s *Stack) hasStack(other *Stack) bool {
	for _, child := range s.children {
		if v, ok := child.(*stackWrapper); ok {
			if v.Stack == other || v.hasStack(other) {
				return true
			}
		}
	}
	return false
}
