package app

import (
	"image"
	"image/color"

	"github.com/gchumillas/prism/ui"
)

type AlertDialog struct {
	app  *App
	text string

	content  *ui.Stack
	closeBtn *ui.Button
	closeId  int

	panel *ui.Stack

	backdrop    *ui.Stack
	backdropBtn *ui.Button
	backdropId  int
}

func (a *App) NewAlertDialog(text string) *AlertDialog {
	ad := &AlertDialog{app: a, text: text}
	ad.content = ui.NewVStack()
	ad.panel = ui.NewVStack()
	ad.backdrop = ui.NewVStack()
	return ad
}

func (ad *AlertDialog) Open() {
	// build backdrop
	b := ad.app.Canvas.Image().Bounds()
	ad.backdropBtn = ui.NewButton(ui.NewRectangle(ad.app.Canvas, b.Dx(), b.Dy(), color.Transparent))
	ad.backdropId = ad.app.Touch.AddTapHandler(ad.backdropBtn, func() {})
	_ = ad.backdropBtn.AttachTo(ad.backdrop)

	// build content
	ad.content.Add(ui.NewText(ad.app.Canvas, ad.text, 600))
	ad.closeBtn = NewLabelButton(ad.app.Canvas, "Close")
	ad.closeId = ad.app.Touch.AddTapHandler(ad.closeBtn, func() {
		ad.Close()
		_ = ad.app.Refresh()
	})
	_ = ad.closeBtn.AttachTo(ad.content)

	// build panel
	size := ad.content.Size()
	ad.panel.Add(ui.NewRectangle(ad.app.Canvas, size.X, size.Y, color.Gray{Y: 200}))
}

func (ad *AlertDialog) Close() {
	// clear backdrop
	ad.backdrop.Clear()
	ad.app.Touch.RemoveTapHandler(ad.backdropId)

	// clear panel
	ad.panel.Clear()

	// clear content
	ad.content.Clear()
	ad.app.Touch.RemoveTapHandler(ad.closeId)
}

func (ad *AlertDialog) Render() []image.Rectangle {
	var dirtyRegs []image.Rectangle

	// render backdrop
	ad.backdrop.Render(0, 0)

	// render background
	regs := ad.panel.Render(500, 500)
	dirtyRegs = append(dirtyRegs, regs...)

	// render content
	regs = ad.content.Render(500, 500)
	dirtyRegs = append(dirtyRegs, regs...)

	return ui.CombineRegions(dirtyRegs)
}
