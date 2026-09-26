package main

import (
	"fmt"
	"image"
	"os"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/internal/platform"
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
	t := app.Region{X0: b.Min.X, Y0: b.Min.Y, X1: b.Max.X, Y1: b.Max.Y}
	return prism.Screen.Refresh([]app.Region{t})
}

func main() {
	env, err := platform.LoadRuntime()
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
