package app

import (
	"image"

	"github.com/gchumillas/prism/ui"
)

type (
	WithFillColor   uint32
	WithStrokeColor uint32
	WithLineWidth   float64
)

const (
	Black       = uint32(0x000000_ff)
	White       = uint32(0xffffff_ff)
	Transparent = uint32(0x000000_00)
)

type Renderer interface {
	Refresh([]image.Rectangle) error
}

type Tappable interface {
	Contains(x, y int) bool
}

type TapDispatcher interface {
	AddTapHandler(btn Tappable, action func()) (id int)
	RemoveTapHandler(id int) error
	Listen()
}

type App struct {
	Touch    TapDispatcher
	Screen   Renderer
	Canvas   *image.RGBA
	FontSize float64

	vs *ui.Stack
	ad *Modal
}

func NewApp(trackpad TapDispatcher, screen Renderer, canvas *image.RGBA, fontSize float64) *App {
	a := &App{
		Touch:    trackpad,
		Screen:   screen,
		Canvas:   canvas,
		FontSize: fontSize,
	}

	a.vs = ui.NewVStack(ui.WithPaddingX(100), ui.WithPaddingY(50))
	btn := NewLabelButton(canvas, "Click me!", ui.WithFontSize(a.FontSize))
	_ = btn.AttachTo(a.vs)

	a.ad = a.NewAlertDialog("Hi there!")
	a.Touch.AddTapHandler(btn, func() {
		a.ad.Open()
		_ = a.Refresh()
	})

	return a
}

func (a *App) Refresh() error {
	// clear the canvas
	b := a.Canvas.Bounds()
	bg := NewRectangle(a.Canvas, b.Dx(), b.Dy(), WithFillColor(White))
	bg.Render(0, 0)

	var dirtyRegs []image.Rectangle
	regs := a.vs.Render(50, 100)
	dirtyRegs = append(dirtyRegs, regs...)

	regs = a.ad.Render()
	dirtyRegs = append(dirtyRegs, regs...)

	dirtyRegs = ui.CombineRegions(dirtyRegs)
	return a.Screen.Refresh(dirtyRegs)
}
