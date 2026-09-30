package test

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/gchumillas/prism/ui"
)

type App struct {
	ui.BaseGroup
	ui.Context
	screen *image.RGBA
}

func NewApp() *App {
	width, height := 640, 480
	ctx := ui.NewContext(width, height, ui.FontSize(10))
	app := &App{
		Context: ctx,
		screen:  image.NewRGBA(image.Rect(0, 0, width, height)),
	}
	app.Add(app.Rectangle(width, height, color.White))
	return app
}

func (app *App) Screen() *image.RGBA {
	return app.screen
}

func (app *App) render() (dirtyRegs []image.Rectangle) {
	for _, c := range app.Children() {
		dirtyRegs = append(dirtyRegs, c.Render()...)
	}
	return ui.MergeRegions(dirtyRegs)
}

func (app *App) Refresh() {
	img := app.Image()
	dirtyRegs := app.render()
	for _, r := range dirtyRegs {
		draw.Draw(app.screen, r, img, r.Min, draw.Src)
	}
}

func (app *App) RefreshRect(r image.Rectangle) {
	app.render()
	img := app.Image()
	draw.Draw(app.screen, r, img, r.Min, draw.Src)
}
