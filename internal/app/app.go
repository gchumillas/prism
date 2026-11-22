package app

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gchumillas/prism/ui"
)

// TODO: rename to Displayer
type Renderer interface {
	Refresh([]image.Rectangle) error
}

type Tappable interface {
	Contains(x, y int) bool
}

type TapDispatcher interface {
	AddTapHandler(btn Tappable, action func()) (id int)
	RemoveTapHandler(id int)
	Listen()
}

type App struct {
	Touch    TapDispatcher
	Screen   Renderer
	Canvas   *image.RGBA
	FontSize float64

	bg *ui.Rectangle
	vs *ui.Stack
}

func NewApp(trackpad TapDispatcher, screen Renderer, canvas *image.RGBA, fontSize float64) *App {
	a := &App{
		Touch:    trackpad,
		Screen:   screen,
		Canvas:   canvas,
		FontSize: fontSize,
	}

	c := ui.NewCanvasFromRGBA(a.Canvas)
	b := a.Canvas.Bounds()
	a.bg = ui.NewRectangle(c, b.Dx(), b.Dy(), color.White)

	a.vs = ui.NewVStack()
	btn := NewLabelButton(c, "Click me!", ui.WithFontSize(a.FontSize))
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
