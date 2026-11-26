package main

import (
	"fmt"
	"image"
	"os"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/internal/platform"
	"github.com/gchumillas/prism/ui"
)

type Storybook struct {
	canvas *image.RGBA
	vs     *ui.Stack
	// physical devices
	Touch    app.TapDispatcher
	Screen   app.Renderer
	FontSize float64
}

func (sb *Storybook) Run() error {
	sb.vs = ui.NewVStack()
	btn := app.NewLabelButton(sb.canvas, "Click me!", ui.WithFontSize(sb.FontSize))
	_ = btn.AttachTo(sb.vs)

	sb.Touch.AddTapHandler(btn, func() {
		fmt.Println("clicked!")
	})

	// refresh the screen
	return sb.refresh()
}

func (sb *Storybook) Render() {
	r := sb.canvas.Bounds()
	app.NewRectangle(sb.canvas, r.Dx(), r.Dy(), app.WithFillColor(app.White))
	sb.vs.Render(100, 200)
}

func (sb *Storybook) refresh() error {
	sb.Render()
	b := sb.canvas.Bounds()
	return sb.Screen.Refresh([]image.Rectangle{b})
}

func main() {
	env, err := platform.LoadRuntime()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	sb := &Storybook{
		canvas:   env.Canvas,
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
