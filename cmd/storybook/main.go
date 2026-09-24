package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/internal/platform"

	"github.com/gchumillas/prism/internal/app"
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
	s := ui.NewScene(ctx, ui.FontSize(sb.FontSize))

	btn := s.LabelButton("Click me!")
	btn.SetX(100)
	btn.SetY(100)
	sb.Add(btn)

	sb.Touch.AddTapHandler(btn, func() {
		fmt.Println("click!")
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
	r := ui.Rectangle{
		Min: ui.Point{X: float64(b.Min.X), Y: float64(b.Min.Y)},
		Max: ui.Point{X: float64(b.Max.X), Y: float64(b.Max.Y)},
	}
	s := app.Region{X0: r.Min.X, Y0: r.Min.Y, X1: r.Max.X, Y1: r.Max.Y}
	return sb.Screen.Refresh([]app.Region{s})
}

func main() {
	env, err := platform.LoadRuntimeEnv()
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
