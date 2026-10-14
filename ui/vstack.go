package ui

import (
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

type stackOption interface {
	applyStack(*VStack)
}

func (val WithPadding) applyStack(vs *VStack) { vs.padding = image.Point(val) }
func (val WithGap) applyStack(vs *VStack)     { vs.gap = int(val) }

type VStack struct {
	ctx             Context
	padding         image.Point
	gap             int
	children        []Widget
	prevChildStates []childState
}

func NewVStack(ctx Context, opts ...stackOption) *VStack {
	vs := &VStack{ctx: ctx}
	for _, opt := range opts {
		opt.applyStack(vs)
	}
	return vs
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

func (vs *VStack) Add(child Widget) {
	vs.children = append(vs.children, child)
}

func (vs *VStack) Remove(child Widget) {
	vs.children = slices.DeleteFunc(vs.children, func(c Widget) bool {
		return c == child
	})
}

func (vs *VStack) AddButton(child Widget) *Button {
	btn := newClickableWidget(child)
	vs.children = append(vs.children, btn)
	return btn.Button
}

func (vs *VStack) AddImageButton(r io.Reader) (*Button, error) {
	img, err := loadImageButton(vs.ctx, r)
	if err != nil {
		return nil, err
	}
	btn := newClickableWidget(img)
	vs.children = append(vs.children, btn)
	return btn.Button, nil
}

func (vs *VStack) AddLabelButton(text string) *Button {
	btn := newClickableWidget(newLabelButton(vs.ctx, text))
	vs.children = append(vs.children, btn)
	return btn.Button
}

func (vs *VStack) RemoveButton(btn *Button) {
	i := slices.IndexFunc(vs.children, func(x Widget) bool {
		if v, ok := x.(*clickableWidget); ok {
			if v.Button == btn {
				return true
			}
		}
		return false
	})
	if i > -1 {
		vs.children = slices.Delete(vs.children, i, i+1)
	}
}

func (vs *VStack) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	pos := image.Pt(x, y)
	childPos := image.Pt(pos.X+vs.padding.X, pos.Y+vs.padding.Y)
	prevHeight := 0
	for i, child := range vs.children {
		if i > 0 {
			childPos.Y += vs.gap + prevHeight
			size.Y += vs.gap
		}
		regs, childSize := child.Render(childPos.X, childPos.Y)
		prevHeight = childSize.Y
		size.X = max(size.X, childSize.X)
		size.Y += childSize.Y
		if i < len(vs.prevChildStates) {
			prev := vs.prevChildStates[i]
			if prev.pos == childPos && prev.typ == reflect.TypeOf(child) && slices.Equal(prev.footprint, child.Footprint()) {
				continue
			}
		}

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
