package storybook

import (
	"fmt"
	"image/draw"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/ui"
)

type Storybook struct {
	Touch    app.TapDispatcher
	Screen   app.Refresher
	Canvas   draw.Image
	FontSize float64
	DPI      float64
}

func (sb *Storybook) Run() error {
	sc := ui.NewScene(sb.Canvas)
	sc.SetDefaultFontSize(sb.FontSize)
	sc.SetDPI(sb.DPI)

	btn := sc.LabelButton("Small")
	btn.SetPosition(50, 100)
	sb.Touch.AddTapHandler(btn, func() {
		fmt.Println("tapped!")
	})
	btn.Draw()

	// refresh the screen
	b := sb.Canvas.Bounds()
	r := app.Region{X0: b.Min.X, Y0: b.Min.Y, X1: b.Max.X, Y1: b.Max.Y}
	return sb.Screen.Refresh([]app.Region{r})
}
