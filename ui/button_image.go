package ui

import (
	"image"
	"image/color"
	"io"

	"github.com/fogleman/gg"
)

var imgBtnPad = image.Pt(5, 5)

type ImageButton struct {
	canvas *gg.Context
	hitbox image.Rectangle
	vstack *VStack
	img    *Image
}

func (c Context) ImageButtonFromReader(r io.Reader) (_ *ImageButton, err error) {
	ib := &ImageButton{canvas: c.canvas}
	ib.vstack = c.VStack(Px(imgBtnPad.X), Py(imgBtnPad.Y))
	if ib.img, err = c.ImageFromReader(r); err != nil {
		return
	}
	ib.vstack.Add(ib.img)
	return ib, nil
}

func (ib *ImageButton) X() int       { return ib.vstack.X() }
func (ib *ImageButton) SetX(val int) { ib.vstack.SetX(val) }

func (ib *ImageButton) Y() int       { return ib.vstack.Y() }
func (ib *ImageButton) SetY(val int) { ib.vstack.SetY(val) }

func (ib *ImageButton) Height() int {
	return ib.img.Height() + 2*imgBtnPad.Y
}

func (ib *ImageButton) Render() image.Rectangle {
	pos := ib.vstack.position
	size := ib.img.Size().Add(imgBtnPad.Mul(2))

	ib.hitbox = image.Rectangle{
		Min: pos,
		Max: pos.Add(size),
	}

	lineWidth := 4
	ib.canvas.SetLineWidth(float64(lineWidth))
	ib.canvas.SetColor(color.Black)
	ib.canvas.DrawRectangle(
		float64(ib.hitbox.Min.X), float64(ib.hitbox.Min.Y),
		float64(ib.hitbox.Dx()), float64(ib.hitbox.Dy()))
	ib.canvas.Stroke()

	r := ib.vstack.Render().Union(ib.hitbox)
	margin := image.Pt(lineWidth/2+1, lineWidth/2+1)
	return image.Rectangle{
		Min: r.Min.Sub(margin),
		Max: r.Max.Add(margin),
	}
}

func (ib *ImageButton) Contains(x, y int) bool {
	if ib.hitbox.Empty() {
		panic("The hitbox is empty. Did you forget calling Render?")
	}
	p := image.Pt(x, y)
	return p.In(ib.hitbox)
}
