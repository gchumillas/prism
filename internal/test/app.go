package test

import (
	"image"
	"image/color"
	"image/draw"
	"slices"

	"github.com/gchumillas/prism/ui"
)

type App struct {
	ui.Context
	screen   *image.RGBA
	children []ui.Renderer
}

func NewApp() *App {
	width, height := 640, 480
	ctx := ui.NewContext(width, height)
	ctx.FontSize = 10
	app := &App{
		Context: ctx,
		screen:  image.NewRGBA(image.Rect(0, 0, width, height)),
	}
	app.Add(ui.NewRectangle(app.Context, width, height, color.White))
	return app
}

func (app *App) Screen() *image.RGBA {
	return app.screen
}

func (app *App) Add(child ui.Renderer) {
	app.children = append(app.children, child)
}

func (app *App) Remove(child ui.Renderer) {
	app.children = slices.DeleteFunc(app.children, func(c ui.Renderer) bool {
		return c == child
	})
}

func (app *App) render() (dirtyRegs []image.Rectangle) {
	for _, c := range app.children {
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
