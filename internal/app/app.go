package app

import (
	"fmt"
	"image/draw"

	"github.com/gchumillas/prism/ui"
)

type Region struct {
	X0, Y0, X1, Y1 int
}

// TODO: rename to Displayer
type Renderer interface {
	Refresh([]Region) error
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
	Canvas   draw.Image
	FontSize float64
	DPI      float64
}

func (app *App) Run() error {
	sc := ui.NewScene(app.Canvas)
	sc.SetDefaultFontSize(app.FontSize)
	sc.SetDPI(app.DPI)

	btn := sc.LabelButton("Click me!")
	btn.SetPosition(50, 100)
	app.Touch.AddTapHandler(btn, func() {
		fmt.Println("tapped!")
	})
	btn.Draw()

	// refresh the screen
	b := app.Canvas.Bounds()
	r := Region{b.Min.X, b.Min.Y, b.Max.X, b.Max.Y}
	return app.Screen.Refresh([]Region{r})
}
