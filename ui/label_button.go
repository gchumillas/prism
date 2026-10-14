package ui

import (
	"hash/maphash"
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

var labelBtnPad = image.Point{38, 19}

type LabelButton struct {
	canvas *gg.Context
	label  *Label
	pos    image.Point
	size   image.Point
}

func NewLabelButton(c Context, text string) *LabelButton {
	lb := &LabelButton{canvas: c.canvas}
	lb.label = NewLabel(c, text)
	lb.size = lb.label.size().Add(labelBtnPad.Mul(2))
	return lb
}

func (lb *LabelButton) Footprint() []uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, lb.pos)
	return append([]uint64{h.Sum64()}, lb.label.Footprint()...)
}

func (lb *LabelButton) Render(x, y int) (dirtyRegs []image.Rectangle, size image.Point) {
	lb.pos = image.Pt(x, y)
	hb := lb.hitbox()
	lineWidth := 4
	lb.canvas.SetLineWidth(float64(lineWidth))
	lb.canvas.DrawRectangle(
		float64(hb.Min.X), float64(hb.Min.Y),
		float64(hb.Dx()), float64(hb.Dy()))
	lb.canvas.SetColor(color.Black)
	lb.canvas.Stroke()

	// addSafeMargin(r, margin)
	margin := image.Pt(lineWidth/2+1, lineWidth/2+1)
	// TODO: use an Outset(rect, margin) utility instead (affect others)
	dirtyReg := image.Rectangle{
		Min: hb.Min.Sub(margin),
		Max: hb.Max.Add(margin),
	}

	dr, _ := lb.label.Render(x+labelBtnPad.X, y+labelBtnPad.Y)
	return append(dr, dirtyReg), image.Pt(hb.Dx(), hb.Dy())
}

func (lb *LabelButton) Contains(x, y int) bool {
	return image.Pt(x, y).In(lb.hitbox())
}

func (lb *LabelButton) hitbox() image.Rectangle {
	return image.Rectangle{
		Min: lb.pos,
		Max: lb.pos.Add(lb.size),
	}
}
