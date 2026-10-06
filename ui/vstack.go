package ui

import (
	"bytes"
	"encoding/json"
	"image"
	"maps"
	"slices"
)

type childState struct {
	y     int
	state map[string]any
}

type VStack struct {
	Padding image.Point
	Gap     int

	children []Widget

	prevPos         image.Point
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

func (vs *VStack) State() map[string]any {
	var childStates []map[string]any
	for _, child := range vs.children {
		childStates = append(childStates, child.State())
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

func (vs *VStack) Add(child Widget) {
	vs.children = append(vs.children, child)
}

// TODO: This changes Remove from deleting every matching occurrence (the previous
// slices.DeleteFunc behavior) to deleting only the first. Since Add permits the
// same widget more than once, Remove(widget) now unexpectedly leaves duplicates
// behind; delete all matches while keeping the cache aligned.
func (vs *VStack) Remove(child Widget) {
	if i := slices.Index(vs.children, child); i > -1 {
		vs.children = append(vs.children[:i], vs.children[i+1:]...)
		if i < len(vs.prevChildStates) {
			vs.prevChildStates = append(vs.prevChildStates[:i], vs.prevChildStates[i+1:]...)
		}
	}
}

func (vs *VStack) Render(x, y int) (dirtyRegs []image.Rectangle) {
	pos := image.Pt(x, y)
	h := pos.Y + vs.Padding.Y
	for i, child := range vs.children {
		if i > 0 {
			h += vs.Gap + vs.children[i-1].Size().Y
		}
		state := child.State()
		prev, ok := childState{}, i < len(vs.prevChildStates)
		if ok {
			prev = vs.prevChildStates[i]
			vs.prevChildStates[i] = childState{y: h, state: state}
		} else {
			vs.prevChildStates = append(vs.prevChildStates, childState{y: h, state: state})
		}
		if vs.prevPos != pos || !ok || prev.y != h || !maps.Equal(prev.state, state) {
			dirtyRegs = append(dirtyRegs, child.Render(pos.X+vs.Padding.X, h)...)
		}
	}
	vs.prevPos = pos
	return
}
