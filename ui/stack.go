package ui

import (
	"hash/maphash"
	"image"
	"io"
	"reflect"
	"slices"
)

type childState struct {
	pos       image.Point
	regs      []image.Rectangle
	typ       reflect.Type
	footprint []uint64
}

type stack struct {
	ctx             Context
	padding         image.Point
	gap             int
	children        []Widget
	prevChildStates []childState
}

type stackOption interface {
	applyStack(*stack)
}

func (val WithPadding) applyStack(s *stack) { s.padding = image.Point(val) }
func (val WithGap) applyStack(s *stack)     { s.gap = int(val) }

func (s *stack) Children() []Widget {
	return s.children
}

func (s *stack) Add(child Widget) {
	s.children = append(s.children, child)
}

func (s *stack) Remove(child Widget) {
	s.children = slices.DeleteFunc(s.children, func(c Widget) bool {
		return c == child
	})
}

func (s *stack) AddButton(child Widget) *Button {
	btn := newClickableWidget(child)
	s.children = append(s.children, btn)
	return btn.Button
}

func (s *stack) AddImageButton(r io.Reader) (*Button, error) {
	img, err := loadImageButton(s.ctx, r)
	if err != nil {
		return nil, err
	}
	btn := newClickableWidget(img)
	s.children = append(s.children, btn)
	return btn.Button, nil
}

func (s *stack) AddLabelButton(text string) *Button {
	btn := newClickableWidget(newLabelButton(s.ctx, text))
	s.children = append(s.children, btn)
	return btn.Button
}

func (s *stack) RemoveButton(btn *Button) {
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
	}
}

func (s *stack) Footprint() []uint64 {
	var result []uint64
	for _, child := range s.children {
		typ := typeID(reflect.TypeOf(child))
		footprint := child.Footprint()
		result = append(result, uint64(typ), uint64(len(footprint)))
		result = append(result, footprint...)
	}

	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, s.padding.X)
	maphash.WriteComparable(&h, s.padding.Y)
	maphash.WriteComparable(&h, s.gap)
	return append([]uint64{h.Sum64()}, result...)
}
