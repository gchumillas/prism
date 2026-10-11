package main

import (
	"fmt"
	"image/color"
	"os"
	"slices"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/internal/platform"
	"github.com/gchumillas/prism/ui"
)

type Storybook struct {
	ctx      ui.Context
	children []ui.Widget
	// physical devices
	Touch  app.TapDispatcher
	Screen app.Renderer
}

func (sb *Storybook) Run() error {
	btn := ui.NewLabelButton(sb.ctx, "Click me!")
	sb.Add(btn)

	sb.Touch.AddTapHandler(btn, func() {
		fmt.Println("clicked!")
	})

	// refresh the screen
	return sb.refresh()
}

func (sb *Storybook) Add(child ui.Widget) {
	sb.children = append(sb.children, child)
}

func (app *Storybook) Remove(child ui.Widget) {
	app.children = slices.DeleteFunc(app.children, func(c ui.Widget) bool {
		return c == child
	})
}

func (sb *Storybook) Render() {
	r := sb.ctx.Canvas().Bounds()
	ui.NewRectangle(sb.ctx, r.Dx(), r.Dy(), color.White)
	for _, e := range sb.children {
		e.Render(100, 200)
	}
}

func (sb *Storybook) refresh() error {
	sb.Render()
	b := sb.ctx.Canvas().Bounds()
	s := app.Region{X0: b.Min.X, Y0: b.Min.Y, X1: b.Max.X, Y1: b.Max.Y}
	return sb.Screen.Refresh([]app.Region{s})
}

func main() {
	env, err := platform.LoadRuntime()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	ctx := ui.NewContextFromRGBA(env.Canvas)
	ctx.FontSize = env.FontSize
	sb := &Storybook{
		ctx:    ctx,
		Touch:  env.Touch,
		Screen: env.Screen,
	}

	if err := sb.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	env.Touch.Listen()
}
