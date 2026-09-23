package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/gchumillas/prism/internal/platform"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/ui"
)

type Storybook struct {
	ui.BaseGroup
	Touch    app.TapDispatcher
	Screen   app.Renderer
	Canvas   draw.Image
	FontSize float64
	DPI      float64
}

func (sb *Storybook) Run() error {
	s := ui.NewScene(sb.Canvas, ui.FontSize(sb.FontSize), ui.DPI(sb.DPI))

	text := "Lorem ipsum dolor lal kfh askd lajs lfkjah sdkj lkjsdff lkasd"
	txt := s.Text(text, 225, ui.FontSize(18))
	txt.SetX(100)
	txt.SetY(100)
	sb.Add(txt)

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
	r := app.Region{X0: b.Min.X, Y0: b.Min.Y, X1: b.Max.X, Y1: b.Max.Y}
	return sb.Screen.Refresh([]app.Region{r})
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
		DPI:      env.DPI,
	}
	if err := sb.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	env.Touch.Listen()
}
