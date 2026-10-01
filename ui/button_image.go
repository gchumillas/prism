package ui

import (
	"image"
	"image/color"
	"io"

	"github.com/fogleman/gg"
)

var imgBtnPad = image.Pt(5, 5)

type ImageButton struct {
	Position image.Point
	Padding  image.Point
	canvas   *gg.Context
	hitbox   image.Rectangle
	vstack   *VStack
	img      *Image
}

func (c Context) ImageButtonFromReader(r io.Reader) (_ *ImageButton, err error) {
	ib := &ImageButton{canvas: c.canvas}
	ib.vstack = c.VStack()
	ib.Padding = imgBtnPad
	if ib.img, err = c.ImageFromReader(r); err != nil {
		return
	}
	ib.vstack.Add(ib.img)
	return ib, nil
}

func (lb *ImageButton) SetX(val int) { lb.Position.X = val }
func (lb *ImageButton) SetY(val int) { lb.Position.Y = val }

func (ib *ImageButton) Height() int {
	return ib.img.Height() + 2*imgBtnPad.Y
}

func (ib *ImageButton) Render() []image.Rectangle {
	ib.vstack.Padding = ib.Padding
	ib.vstack.Position = ib.Position
	size := ib.img.Size().Add(imgBtnPad.Mul(2))

	ib.hitbox = image.Rectangle{
		Min: ib.Position,
		Max: ib.Position.Add(size),
	}

	lineWidth := 4
	ib.canvas.SetLineWidth(float64(lineWidth))
	ib.canvas.SetColor(color.Black)
	ib.canvas.DrawRectangle(
		float64(ib.hitbox.Min.X), float64(ib.hitbox.Min.Y),
		float64(ib.hitbox.Dx()), float64(ib.hitbox.Dy()))
	ib.canvas.Stroke()

	margin := image.Pt(lineWidth/2+1, lineWidth/2+1)
	return append(ib.vstack.Render(), image.Rectangle{
		Min: ib.hitbox.Min.Sub(margin),
		Max: ib.hitbox.Max.Add(margin),
	})
}

func (ib *ImageButton) Contains(x, y int) bool {
	if ib.hitbox.Empty() {
		panic("The hitbox is empty. Did you forget calling Render?")
	}
	p := image.Pt(x, y)
	return p.In(ib.hitbox)
}
