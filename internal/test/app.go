package test

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/gchumillas/prism/ui"
)

type App struct {
	ui.Context
	*ui.VStack
	bg     *ui.Rectangle
	screen *image.RGBA
}

func NewApp() *App {
	width, height := 640, 480
	ctx := ui.NewContext(width, height)
	ctx.FontSize = 10
	app := &App{
		Context: ctx,
		VStack:  ui.NewVStack(ctx),
		bg:      ui.NewRectangle(ctx, width, height, color.White),
		screen:  image.NewRGBA(image.Rect(0, 0, width, height)),
	}
	return app
}

func (app *App) Render() []image.Rectangle {
	p := image.Pt(0, 0)
	return append(app.bg.Render(p), app.VStack.Render(p)...)
}

func (app *App) Screen() *image.RGBA {
	return app.screen
}

func (app *App) Refresh() {
	img := app.Canvas()
	dirtyRegs := app.Render()
	for _, r := range dirtyRegs {
		draw.Draw(app.screen, r, img, r.Min, draw.Src)
	}
}

func (app *App) RefreshRect(r image.Rectangle) {
	app.Render()
	img := app.Canvas()
	draw.Draw(app.screen, r, img, r.Min, draw.Src)
}
