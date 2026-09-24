package main

import (
	"fmt"
	"image"
	"os"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/internal/platform"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/ui"
)

type Prism struct {
	Touch    app.TapDispatcher
	Screen   app.Renderer
	Canvas   *image.RGBA
	FontSize float64
}

func (prism *Prism) Run() error {
	ctx := gg.NewContextForRGBA(prism.Canvas)
	s := ui.NewScene(ctx, ui.FontSize(prism.FontSize))

	btn := s.LabelButton("Click me!")
	btn.SetX(50)
	btn.SetY(100)
	prism.Touch.AddTapHandler(btn, func() {
		fmt.Println("tapped!")
	})
	btn.Render()

	// refresh the screen
	b := prism.Canvas.Bounds()
	r := ui.Rectangle{
		Min: ui.Point{X: float64(b.Min.X), Y: float64(b.Min.Y)},
		Max: ui.Point{X: float64(b.Max.X), Y: float64(b.Max.Y)},
	}
	t := app.Region{X0: r.Min.X, Y0: r.Min.Y, X1: r.Max.X, Y1: r.Max.Y}
	return prism.Screen.Refresh([]app.Region{t})
}

func main() {
	env, err := platform.LoadRuntimeEnv()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	prism := &Prism{
		Touch:    env.Touch,
		Screen:   env.Screen,
		Canvas:   env.Canvas,
		FontSize: env.FontSize,
	}
	if err := prism.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	env.Touch.Listen()
}
