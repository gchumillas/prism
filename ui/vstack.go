package ui

import (
	"bytes"
	"encoding/json"
	"image"
	"maps"
	"reflect"
	"slices"
)

type childState struct {
	pos       image.Point
	regs      []image.Rectangle
	typ       reflect.Type
	footprint map[string]any
}

type VStack struct {
	Padding image.Point
	Gap     int

	children []Widget

	prevChildStates []childState
}

func NewVStack(c Context) *VStack {
	return &VStack{}
}

func (vs *VStack) Size() image.Point {
	width, height := 0, vs.Padding.Y
	for i, child := range vs.children {
		if i > 0 {
			height += vs.Gap
		}
		size := child.Size()
		width = max(width, size.X)
		height += size.Y
	}
	return image.Pt(width+vs.Padding.X, height)
}

func (vs *VStack) Footprint() map[string]any {
	var childStates []struct {
		Typ       int
		Footprint map[string]any
	}
	for _, child := range vs.children {
		childStates = append(childStates, struct {
			Typ       int
			Footprint map[string]any
		}{
			Typ:       typeID(reflect.TypeOf(child)),
			Footprint: child.Footprint(),
		})
	}

	// TODO: this is a temporary solution. We will study how to avoid using the `json` library in a separate PR.
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(childStates); err != nil {
		panic(err)
	}

	return map[string]any{"gap": vs.Gap, "padding": vs.Padding, "children": buff.String()}
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
	childPos := image.Pt(pos.X+vs.Padding.X, pos.Y+vs.Padding.Y)
	for i, child := range vs.children {
		if i > 0 {
			childPos.Y += vs.Gap + vs.children[i-1].Size().Y
		}

		if i < len(vs.prevChildStates) {
			prev := vs.prevChildStates[i]
			if prev.pos == childPos && prev.typ == reflect.TypeOf(child) && maps.Equal(prev.footprint, child.Footprint()) {
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
