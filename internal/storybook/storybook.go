package storybook

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/ui"
)

type Storybook struct {
	ui.Group
	Touch    app.TapDispatcher
	Screen   app.Renderer
	Canvas   draw.Image
	FontSize float64
	DPI      float64
}

func (sb *Storybook) Run() error {
	sc := ui.NewScene(sb.Canvas)
	sc.SetDefaultFontSize(sb.FontSize)
	sc.SetDPI(sb.DPI)

	vl := sc.VStack()
	vl.SetGap(10)
	vl.Add(sc.LabelButton("Option 1"))
	vl.Add(sc.LabelButton("Option 2"))
	vl.Add(sc.LabelButton("Option 3"))
	sb.Add(vl)

	// refresh the screen
	sb.Draw()
	b := sb.Canvas.Bounds()
	r := app.Region{X0: b.Min.X, Y0: b.Min.Y, X1: b.Max.X, Y1: b.Max.Y}
	return sb.Screen.Refresh([]app.Region{r})
}

func (sb *Storybook) Draw() {
	r := sb.Canvas.Bounds()
	draw.Draw(sb.Canvas, r, image.NewUniform(color.White), image.Pt(0, 0), draw.Src)
	for _, e := range sb.Elems {
		e.Draw()
	}
}
