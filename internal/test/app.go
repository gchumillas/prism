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
	width, height := 320, 240
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

func (app *App) render() {
	for _, c := range app.Children() {
		c.Render()
	}
}

func (app *App) Refresh() {
	app.render()
	img := app.Image()
	draw.Draw(app.screen, img.Bounds(), img, image.Pt(0, 0), draw.Src)
}

func (app *App) RefreshRect(r image.Rectangle) {
	app.render()
	img := app.Image()
	draw.Draw(app.screen, r, img, r.Min, draw.Src)
}
