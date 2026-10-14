package ui

import (
	"hash/maphash"
	"image"
	"io"
	"reflect"
	"slices"
)

type HStack struct {
	ctx             Context
	padding         image.Point
	gap             int
	children        []Widget
	prevChildStates []childState
}

type hStackOption interface {
	applyHStack(*HStack)
}

func (val WithPadding) applyHStack(vs *HStack) { vs.padding = image.Point(val) }
func (val WithGap) applyHStack(vs *HStack)     { vs.gap = int(val) }

func NewHStack(ctx Context, opts ...hStackOption) *HStack {
	vs := &HStack{}
	vs.ctx = ctx
	for _, opt := range opts {
		opt.applyHStack(vs)
	}
	return vs
}

func (hs *HStack) Children() []Widget {
	return hs.children
}

func (hs *HStack) Add(child Widget) {
	hs.children = append(hs.children, child)
}

func (hs *HStack) Remove(child Widget) {
	hs.children = slices.DeleteFunc(hs.children, func(c Widget) bool {
		return c == child
	})
}

func (hs *HStack) AddButton(child Widget) *Button {
	btn := newClickableWidget(child)
	hs.children = append(hs.children, btn)
	return btn.Button
}

func (hs *HStack) AddImageButton(r io.Reader) (*Button, error) {
	img, err := loadImageButton(hs.ctx, r)
	if err != nil {
		return nil, err
	}
	btn := newClickableWidget(img)
	hs.children = append(hs.children, btn)
	return btn.Button, nil
}

func (hs *HStack) AddLabelButton(text string) *Button {
	btn := newClickableWidget(newLabelButton(hs.ctx, text))
	hs.children = append(hs.children, btn)
	return btn.Button
}

func (hs *HStack) RemoveButton(btn *Button) {
	i := slices.IndexFunc(hs.children, func(x Widget) bool {
		if v, ok := x.(*clickableWidget); ok {
			if v.Button == btn {
				return true
			}
		}
		return false
	})
	if i > -1 {
		hs.children = slices.Delete(hs.children, i, i+1)
	}
}

func (hs *HStack) Footprint() []uint64 {
	var result []uint64
	for _, child := range hs.children {
		typ := typeID(reflect.TypeOf(child))
		footprint := child.Footprint()
		result = append(result, uint64(typ), uint64(len(footprint)))
		result = append(result, footprint...)
	}

	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, hs.padding.X)
	maphash.WriteComparable(&h, hs.padding.Y)
	maphash.WriteComparable(&h, hs.gap)
	return append([]uint64{h.Sum64()}, result...)
}

func (hs *HStack) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	size = image.Pt(hs.padding.X, 0)
	for i, child := range hs.children {
		if i > 0 {
			size.X += hs.gap
		}
		pos := image.Pt(x+size.X, y+hs.padding.Y)
		childRegs, childSize := child.Render(pos.X, pos.Y)
		cs := childState{
			pos:       pos,
			typ:       reflect.TypeOf(child),
			footprint: child.Footprint(),
			// TODO: this is wrong, since "size" and "dirty regions" are different concepts
			regs: []image.Rectangle{{Min: pos, Max: pos.Add(childSize)}},
		}
		if i < len(hs.prevChildStates) {
			prev := &hs.prevChildStates[i]
			if prev.pos != cs.pos || prev.typ != cs.typ || !slices.Equal(prev.footprint, cs.footprint) {
				dirtyRegs = append(dirtyRegs, prev.regs...)
				dirtyRegs = append(dirtyRegs, childRegs...)
				*prev = cs
			}
		} else {
			hs.prevChildStates = append(hs.prevChildStates, cs)
			dirtyRegs = append(dirtyRegs, childRegs...)
		}
		size.X += childSize.X
		size.Y = max(size.Y, childSize.Y)
	}
	for _, child := range hs.prevChildStates[len(hs.children):] {
		dirtyRegs = append(dirtyRegs, child.regs...)
	}
	// TODO: button leaks?
	hs.prevChildStates = hs.prevChildStates[:len(hs.children)]
	dirtyRegs = MergeRegions(dirtyRegs)
	size = image.Pt(size.X+hs.padding.X, size.Y+2*hs.padding.Y)
	return
}
