package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/internal/platform"
	"github.com/gchumillas/prism/ui"
)

type Storybook struct {
	ui.BaseGroup
	Touch    app.TapDispatcher
	Screen   app.Renderer
	Canvas   *image.RGBA
	FontSize float64
}

func (sb *Storybook) Run() error {
	ctx := gg.NewContextForRGBA(sb.Canvas)
	c := ui.NewContext(ctx, ui.FontSize(sb.FontSize))

	btn := c.LabelButton("Click me!")
	btn.SetX(100)
	btn.SetY(200)
	sb.Add(btn)

	sb.Touch.AddTapHandler(btn, func() {
		fmt.Println("clicked!")
	})

	// refresh the screen
	return sb.refresh()
}

func (sb *Storybook) Render() {
	r := sb.Canvas.Bounds()
	draw.Draw(sb.Canvas, r, image.NewUniform(color.White), image.Pt(0, 0), draw.Src)
	for _, e := range sb.Children() {
		e.Render()
	}
}

func (sb *Storybook) refresh() error {
	sb.Render()
	b := sb.Canvas.Bounds()
	s := app.Region{X0: b.Min.X, Y0: b.Min.Y, X1: b.Max.X, Y1: b.Max.Y}
	return sb.Screen.Refresh([]app.Region{s})
}

func main() {
	env, err := platform.LoadRuntime()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	sb := &Storybook{
		Touch:    env.Touch,
		Screen:   env.Screen,
		Canvas:   env.Canvas,
		FontSize: env.FontSize,
	}
	if err := sb.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	env.Touch.Listen()
}
