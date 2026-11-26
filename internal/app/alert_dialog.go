package app

import (
	"image"

	"github.com/gchumillas/prism/ui"
)

type alertDialogContent struct {
	app     *App
	modal   *Modal
	text    string
	content *ui.Stack
	closeId int
}

func (adc *alertDialogContent) Mount() {
	adc.content.Add(ui.NewText(adc.app.Canvas, adc.text, 600))
	btn := NewLabelButton(adc.app.Canvas, "Close")
	adc.closeId = adc.app.Touch.AddTapHandler(btn, func() {
		adc.modal.Close()
		_ = adc.app.Refresh()
	})
	_ = btn.AttachTo(adc.content)
}

func (adc *alertDialogContent) Unmount() {
	adc.content.Clear()
	adc.app.Touch.RemoveTapHandler(adc.closeId)
}

func (adc *alertDialogContent) Render(x, y int) []image.Rectangle {
	return adc.content.Render(x, y)
}

func (adc *alertDialogContent) Size() image.Point {
	return adc.content.Size()
}

func (a *App) NewAlertDialog(text string) *Modal {
	ad := &alertDialogContent{
		app:     a,
		text:    text,
		content: ui.NewVStack(ui.WithPaddingX(25), ui.WithPaddingY(25), ui.WithGap(25)),
	}
	ad.modal = a.NewModal(ad)
	return ad.modal
}
