package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/internal/platform"
	"github.com/gchumillas/prism/ui"
)

type App struct {
	Touch    app.TapDispatcher
	Screen   app.Renderer
	Canvas   *image.RGBA
	FontSize float64

	bg *ui.Rectangle
	vs *ui.Stack
}

func NewApp(env *platform.Env) *App {
	a := &App{
		Touch:    env.Touch,
		Screen:   env.Screen,
		Canvas:   env.Canvas,
		FontSize: env.FontSize,
	}

	c := ui.NewCanvasFromRGBA(a.Canvas)
	b := a.Canvas.Bounds()
	a.bg = ui.NewRectangle(c, b.Dx(), b.Dy(), color.White)

	a.vs = ui.NewVStack()
	btn := app.NewLabelButton(c, "Click me!", ui.WithFontSize(a.FontSize))
	_ = btn.AttachTo(a.vs)

	a.Touch.AddTapHandler(btn, func() {
		fmt.Println("tapped!")
	})

	return a
}

func (a *App) Refresh() error {
	var dirtyRegs []image.Rectangle

	regs, _ := a.bg.Render(0, 0)
	dirtyRegs = append(dirtyRegs, regs...)

	regs, _ = a.vs.Render(50, 100)
	dirtyRegs = append(dirtyRegs, regs...)

	dirtyRegs = ui.CombineRegions(dirtyRegs)
	return a.Screen.Refresh(dirtyRegs)
}

func main() {
	env, err := platform.LoadRuntime()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	a := NewApp(env)
	if err := a.Refresh(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	env.Touch.Listen()
}
