package ui

import (
	"fmt"
	"image"
	"reflect"
	"slices"
)

type drawable interface {
	Render(x, y int) (dirtyRegs []image.Rectangle)
	Size() image.Point
}

type childState struct {
	pos       image.Point
	regs      []image.Rectangle
	typ       reflect.Type
	footprint []uint64
}

type stackContext struct {
	prevChildStatesIdx int
	prevChildStates    []childState
}

type Stack struct {
	PaddingX        int
	PaddingY        int
	Gap             int
	children        []drawable
	prevChildStates []childState
	layout          func(*stackContext, *Stack, int, int) ([]image.Rectangle, image.Point)
	layoutSize      func(*Stack) image.Point
	outdated        bool
	attached        bool
}

func newStack(layout func(*stackContext, *Stack, int, int) ([]image.Rectangle, image.Point), layoutSize func(*Stack) image.Point) *Stack {
	return &Stack{layout: layout, layoutSize: layoutSize}
}

func (s *Stack) Add(child Widget) {
	s.children = append(s.children, child)
}

func (s *Stack) Remove(child Widget) (n int) {
	oldLen := len(s.children)
	s.children = slices.DeleteFunc(s.children, func(c drawable) bool {
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
	if i := slices.IndexFunc(parent.children, func(x drawable) bool {
		return x == s
	}); i > -1 {
		parent.children = slices.Delete(parent.children, i, i+1)
		s.attached = false
		parent.outdated = true
		return nil
	}
	return ErrNotAttached
}

func (s *Stack) Render(x, y int) (dirtyRegs []image.Rectangle) {
	stx := &stackContext{prevChildStates: s.prevChildStates}
	dirtyRegs, _ = s.layout(stx, s, x, y)

	// discarded regions must be refreshed
	for _, child := range stx.prevChildStates[stx.prevChildStatesIdx:] {
		dirtyRegs = append(dirtyRegs, child.regs...)
	}

	// clear discarded childs
	clear(stx.prevChildStates[stx.prevChildStatesIdx:])
	stx.prevChildStates = stx.prevChildStates[:stx.prevChildStatesIdx]

	dirtyRegs = CombineRegions(dirtyRegs)
	s.prevChildStates = stx.prevChildStates
	return
}

func (s *Stack) Size() image.Point {
	return s.layoutSize(s)
}

func (s *Stack) Clear() {
	for _, child := range slices.Clone(s.children) {
		if v, ok := child.(interface{ DetachFrom(*Stack) error }); ok {
			if err := v.DetachFrom(s); err != nil {
				panic(err)
			}
		}
	}
	clear(s.children)
	s.children = s.children[:0]
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
