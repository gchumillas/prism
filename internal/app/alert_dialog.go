package app

import (
	"image"
	"image/color"

	"github.com/gchumillas/prism/ui"
)

type AlertDialog struct {
	app  *App
	vs   *ui.Stack
	bg   *ui.Stack
	text string

	msg      *ui.Text
	closeBtn *ui.Button
	bgRec    *ui.Rectangle
	closeId  int
}

func (a *App) NewAlertDialog(text string) *AlertDialog {
	ad := &AlertDialog{app: a, text: text}
	ad.vs = ui.NewVStack()
	ad.bg = ui.NewVStack()
	return ad
}

func (ad *AlertDialog) Open() {
	// build interface
	ad.msg = ui.NewText(ad.app.Canvas, ad.text, 600)
	ad.vs.Add(ad.msg)
	ad.closeBtn = NewLabelButton(ad.app.Canvas, "Close")
	ad.closeId = ad.app.Touch.AddTapHandler(ad.closeBtn, func() {
		ad.Close()
		_ = ad.app.Refresh()
	})
	_ = ad.closeBtn.AttachTo(ad.vs)

	// build background
	size := ad.vs.Size()
	ad.bgRec = ui.NewRectangle(ad.app.Canvas, size.X, size.Y, color.Gray{Y: 200})
	ad.bg.Add(ad.bgRec)
}

func (ad *AlertDialog) Close() {
	// clear background
	ad.bg.Remove(ad.bgRec)

	// clear interface
	ad.vs.Remove(ad.msg)
	_ = ad.closeBtn.DetachFrom(ad.vs)
	ad.app.Touch.RemoveTapHandler(ad.closeId)
}

func (ad *AlertDialog) Render() []image.Rectangle {
	var dirtyRegs []image.Rectangle

	// render background
	regs := ad.bg.Render(500, 500)
	dirtyRegs = append(dirtyRegs, regs...)

	// render interface
	regs = ad.vs.Render(500, 500)
	dirtyRegs = append(dirtyRegs, regs...)

	return ui.CombineRegions(dirtyRegs)
}
