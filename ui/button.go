package ui

import (
	"image"
	"slices"
)

type Button struct {
	bounds image.Rectangle
}

func (b *Button) Contains(x, y int) bool {
	return image.Pt(x, y).In(b.bounds)
}

type clickableWidget struct {
	*Button
	widget Widget
}

func newClickableWidget(w Widget) *clickableWidget {
	return &clickableWidget{Button: &Button{}, widget: w}
}

func AttachButton(owner *Stack, w Widget) *Button {
	btn := newClickableWidget(w)
	owner.children = append(owner.children, btn)
	return btn.Button
}

func (b *Button) DetachFrom(owner *Stack) bool {
	i := slices.IndexFunc(owner.children, func(x Widget) bool {
		if v, ok := x.(*clickableWidget); ok {
			if v.Button == b {
				return true
			}
		}
		return false
	})
	if i > -1 {
		owner.children = slices.Delete(owner.children, i, i+1)
		return true
	}
	return false
}

func (cb *clickableWidget) Render(x, y int) ([]image.Rectangle, image.Point) {
	dirtyRegs, size := cb.widget.Render(x, y)
	pos := image.Pt(x, y)
	cb.bounds = image.Rectangle{Min: pos, Max: pos.Add(size)}
	return dirtyRegs, size
}

func (cb *clickableWidget) Footprint() []uint64 {
	return cb.widget.Footprint()
}
