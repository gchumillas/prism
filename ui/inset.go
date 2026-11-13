package ui

import (
	"image"
)

type Inset struct {
	paddingX int
	paddingY int
	target   Widget
}

func NewInset(w Widget, paddingX, paddingY int) *Inset {
	return &Inset{target: w, paddingX: paddingX, paddingY: paddingY}
}

func (b *Inset) Footprint() []uint64 {
	return append([]uint64{uint64(b.paddingX), uint64(b.paddingY)}, b.target.Footprint()...)
}

func (b *Inset) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	dirtyRegs, size = b.target.Render(x+b.paddingX, y+b.paddingY)
	size.X += 2 * b.paddingX
	size.Y += 2 * b.paddingY
	return
}
