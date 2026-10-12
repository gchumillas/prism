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

type stackOption interface {
	applyStack(*VStack)
}

func (val WithPadding) applyStack(vs *VStack) { vs.padding = image.Point(val) }
func (val WithGap) applyStack(vs *VStack)     { vs.gap = int(val) }

type VStack struct {
	padding         image.Point
	gap             int
	children        []Widget
	prevChildStates []childState
}

func NewVStack(c Context, opts ...stackOption) *VStack {
	vs := &VStack{}
	for _, opt := range opts {
		opt.applyStack(vs)
	}
	return vs
}

func (vs *VStack) Size() image.Point {
	width, height := 0, vs.padding.Y
	for i, child := range vs.children {
		if i > 0 {
			height += vs.gap
		}
		size := child.Size()
		width = max(width, size.X)
		height += size.Y
	}
	return image.Pt(width+vs.padding.X, height)
}

func (vs *VStack) Footprint() []uint64 {
	var result []uint64
	for _, child := range vs.children {
		typ := typeID(reflect.TypeOf(child))
		footprint := child.Footprint()
		result = append(result, uint64(typ), uint64(len(footprint)))
		result = append(result, footprint...)
	}
	return result
}

func (vs *VStack) Children() []Widget {
	return vs.children
}

// TODO: the same child can be added more than once
func (vs *VStack) Add(child Widget) {
	vs.children = append(vs.children, child)
}

func (vs *VStack) Remove(child Widget) {
	vs.children = slices.DeleteFunc(vs.children, func(c Widget) bool {
		return c == child
	})
}

func (vs *VStack) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	childPos := image.Pt(pos.X+vs.padding.X, pos.Y+vs.padding.Y)
	for i, child := range vs.children {
		if i > 0 {
			childPos.Y += vs.gap + vs.children[i-1].Size().Y
		}

		if i < len(vs.prevChildStates) {
			prev := vs.prevChildStates[i]
			if prev.pos == childPos && prev.typ == reflect.TypeOf(child) && slices.Equal(prev.footprint, child.Footprint()) {
				continue
			}
		}

		regs := child.Render(childPos.X, childPos.Y)
		dirtyRegs = append(dirtyRegs, regs...)
		cs := childState{pos: childPos, regs: regs, typ: reflect.TypeOf(child), footprint: child.Footprint()}
		if i < len(vs.prevChildStates) {
			dirtyRegs = append(dirtyRegs, vs.prevChildStates[i].regs...)
			vs.prevChildStates[i] = cs
		} else {
			vs.prevChildStates = append(vs.prevChildStates, cs)
		}
	}
	for _, cs := range vs.prevChildStates[len(vs.children):] {
		dirtyRegs = append(dirtyRegs, cs.regs...)
	}
	vs.prevChildStates = vs.prevChildStates[:len(vs.children)]
	dirtyRegs = MergeRegions(dirtyRegs)
	return
}
