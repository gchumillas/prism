package ui

import (
	"fmt"
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

type stackContext struct {
	prevChildStatesIdx int
	prevChildStates    []childState
	dirtyRegs          []image.Rectangle
	size               image.Point
}

type renderer interface {
	Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point)
}

type Stack struct {
	paddingX        int
	paddingY        int
	gap             int
	children        []renderer
	prevChildStates []childState
	layout          func(*stackContext, *Stack, int, int)
	outdated        bool
	attached        bool
}

type StackOption interface {
	ApplyStack(*Stack)
}

func (val WithPaddingX) ApplyStack(s *Stack) { s.paddingX = int(val) }
func (val WithPaddingY) ApplyStack(s *Stack) { s.paddingY = int(val) }
func (val WithGap) ApplyStack(s *Stack)      { s.gap = int(val) }

func newStack(layout func(*stackContext, *Stack, int, int), opts ...StackOption) *Stack {
	s := &Stack{layout: layout}
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
	s.children = slices.DeleteFunc(s.children, func(c renderer) bool {
		return c == child
	})
	return oldLen - len(s.children)
}

func (s *Stack) AttachTo(parent *Stack) error {
	if s.attached {
		return fmt.Errorf("stack: %w", ErrAttached)
	}
	if parent == s || s.hasStack(parent) {
		return fmt.Errorf("stack: %w", ErrCircularRef)
	}
	if parent.outdated {
		return fmt.Errorf("stack: %w", ErrNeedsRender)
	}
	s.reset()
	parent.children = append(parent.children, s)
	s.attached = true
	return nil
}

func (s *Stack) DetachFrom(parent *Stack) error {
	if i := slices.IndexFunc(parent.children, func(x renderer) bool {
		return x == s
	}); i > -1 {
		parent.children = slices.Delete(parent.children, i, i+1)
		s.attached = false
		parent.outdated = true
		return nil
	}
	return ErrNotAttached
}

func (s *Stack) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	stx := &stackContext{prevChildStates: s.prevChildStates}
	s.layout(stx, s, x, y)

	// discarded regions must be refreshed
	for _, child := range stx.prevChildStates[stx.prevChildStatesIdx:] {
		stx.dirtyRegs = append(stx.dirtyRegs, child.regs...)
	}

	// clear discarded childs
	clear(stx.prevChildStates[stx.prevChildStatesIdx:])
	stx.prevChildStates = stx.prevChildStates[:stx.prevChildStatesIdx]

	stx.dirtyRegs = CombineRegions(stx.dirtyRegs)
	s.prevChildStates = stx.prevChildStates
	dirtyRegs = stx.dirtyRegs
	size = stx.size
	return
}

func (s *Stack) hasStack(other *Stack) bool {
	for _, child := range s.children {
		if v, ok := child.(*Stack); ok {
			if v == other || v.hasStack(other) {
				return true
			}
		}
	}
	return false
}

func (s *Stack) reset() {
	s.prevChildStates = nil
	for _, child := range s.children {
		if v, ok := child.(*Stack); ok {
			v.reset()
		}
	}
}
