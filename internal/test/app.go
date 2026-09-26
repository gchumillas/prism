package test

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/ui"
)

type App struct {
	ui.BaseGroup
	ui.Context
	canvas  *gg.Context
	screen  *image.RGBA
	bgColor color.Color
}

func NewApp(width, height int, bgColor color.Color) *App {
	canvas := gg.NewContext(width, height)
	return &App{
		Context: *ui.NewContext(canvas, ui.FontSize(10)),
		canvas:  canvas,
		screen:  image.NewRGBA(image.Rect(0, 0, width, height)),
		bgColor: bgColor,
	}
}

func (app *App) Screen() *image.RGBA {
	return app.screen
}

func (app *App) render() {
	width, height := app.canvas.Width(), app.canvas.Height()
	app.canvas.DrawRectangle(0, 0, float64(width), float64(height))
	app.canvas.SetColor(app.bgColor)
	app.canvas.Fill()

	for _, c := range app.Children() {
		c.Render()
	}
}

func (app *App) Refresh() {
	app.render()
	img := app.canvas.Image()
	draw.Draw(app.screen, img.Bounds(), img, image.Pt(0, 0), draw.Src)
}

func (app *App) RefreshRect(r image.Rectangle) {
	app.render()
	img := app.canvas.Image()
	draw.Draw(app.screen, r, img, r.Min, draw.Src)
}
