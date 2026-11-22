package app

import (
	"image"

	"github.com/gchumillas/prism/ui"
)

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

	vs *ui.Stack
}

func NewApp(trackpad TapDispatcher, screen Renderer, canvas *image.RGBA, fontSize float64) *App {
	c := ui.NewCanvasFromRGBA(canvas)
	a := &App{
		Touch:    trackpad,
		Screen:   screen,
		Canvas:   canvas,
		FontSize: fontSize,
	}

	a.vs = ui.NewVStack(ui.WithPaddingX(100), ui.WithPaddingY(50))
	btn := NewLabelButton(c, "Click me!", ui.WithFontSize(a.FontSize))
	_ = btn.AttachTo(a.vs)

	a.Touch.AddTapHandler(btn, func() {
		_ = a.Refresh()
	})

	return a
}

func (a *App) Refresh() error {
	var dirtyRegs []image.Rectangle

	regs, _ := a.vs.Render(50, 100)
	dirtyRegs = append(dirtyRegs, regs...)

	dirtyRegs = ui.CombineRegions(dirtyRegs)
	return a.Screen.Refresh(dirtyRegs)
}
