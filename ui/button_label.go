package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

var labelBtnPad = image.Point{38, 19}

type LabelButton struct {
	canvas *gg.Context
	hitbox image.Rectangle
	vstack *VStack
	label  *Label
}

type LabelButtonOpt interface {
	applyLabelButton(*LabelButton)
}

// avail options
func (val xOpt) applyLabelButton(x *LabelButton) { x.SetX(int(val)) }
func (val yOpt) applyLabelButton(x *LabelButton) { x.SetY(int(val)) }

func (c Context) LabelButton(text string, opts ...LabelButtonOpt) *LabelButton {
	lb := &LabelButton{canvas: c.canvas}
	lb.vstack = c.VStack(Px(labelBtnPad.X), Py(labelBtnPad.Y))
	for _, opt := range opts {
		opt.applyLabelButton(lb)
	}
	lb.label = c.Label(text)
	lb.vstack.Add(lb.label)
	return lb
}

func (lb *LabelButton) X() int       { return lb.vstack.X() }
func (lb *LabelButton) SetX(val int) { lb.vstack.SetX(val) }

func (lb *LabelButton) Y() int       { return lb.vstack.Y() }
func (lb *LabelButton) SetY(val int) { lb.vstack.SetY(val) }

func (lb *LabelButton) Height() int {
	return lb.label.Height() + 2*labelBtnPad.Y
}

func (lb *LabelButton) Render() []image.Rectangle {
	size := image.Pt(lb.label.Width(), lb.label.Height()).Add(labelBtnPad.Mul(2))
	lb.hitbox = image.Rectangle{
		Min: lb.vstack.position,
		Max: lb.vstack.position.Add(size),
	}

	lineWidth := 4
	lb.canvas.SetLineWidth(float64(lineWidth))
	lb.canvas.DrawRectangle(
		float64(lb.hitbox.Min.X), float64(lb.hitbox.Min.Y),
		float64(lb.hitbox.Dx()), float64(lb.hitbox.Dy()))
	lb.canvas.SetColor(color.Black)
	lb.canvas.Stroke()

	margin := image.Pt(lineWidth/2+1, lineWidth/2+1)
	return append(lb.vstack.Render(), image.Rectangle{
		Min: lb.hitbox.Min.Sub(margin),
		Max: lb.hitbox.Max.Add(margin),
	})
}

func (lb *LabelButton) Contains(x, y int) bool {
	if lb.hitbox.Empty() {
		panic("The hitbox is empty. Did you forget calling Render?")
	}
	p := image.Pt(x, y)
	return p.In(lb.hitbox)
}
