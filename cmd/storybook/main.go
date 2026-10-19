package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/internal/platform"
	"github.com/gchumillas/prism/ui"
)

type Storybook struct {
	canvas *ui.Canvas
	vs     *ui.Stack
	// physical devices
	Touch    app.TapDispatcher
	Screen   app.Renderer
	FontSize float64
}

func (sb *Storybook) Run() error {
	sb.vs = ui.NewVStack(sb.canvas)
	btn := app.NewLabelButton(sb.canvas, "Click me!", ui.WithFontSize(sb.FontSize))
	_ = btn.AttachTo(sb.vs)

	sb.Touch.AddTapHandler(btn, func() {
		fmt.Println("clicked!")
	})

	// refresh the screen
	return sb.refresh()
}

func (sb *Storybook) Render() {
	r := sb.canvas.Image().Bounds()
	ui.NewRectangle(sb.canvas, r.Dx(), r.Dy(), color.White)
	sb.vs.Render(100, 200)
}

func (sb *Storybook) refresh() error {
	sb.Render()
	b := sb.canvas.Image().Bounds()
	s := app.Region{X0: b.Min.X, Y0: b.Min.Y, X1: b.Max.X, Y1: b.Max.Y}
	return sb.Screen.Refresh([]app.Region{s})
}

func main() {
	env, err := platform.LoadRuntime()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	c := ui.NewCanvasFromRGBA(env.Canvas)
	sb := &Storybook{
		canvas:   c,
		Touch:    env.Touch,
		Screen:   env.Screen,
		FontSize: env.FontSize,
	}

	if err := sb.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	env.Touch.Listen()
}
