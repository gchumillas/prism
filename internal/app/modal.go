package app

import (
	"image"
	"image/color"

	"github.com/gchumillas/prism/ui"
)

type ModalContent interface {
	ui.Drawable
	Mount()
	Unmount()
}

type Modal struct {
	app        *App
	backdrop   *ui.Stack
	panel      *ui.Stack
	content    ModalContent
	backdropId int
}

func (a *App) NewModal(content ModalContent) *Modal {
	return &Modal{
		app:      a,
		backdrop: ui.NewVStack(),
		panel:    ui.NewVStack(),
		content:  content,
	}
}

func (m *Modal) Open() {
	b := m.app.Canvas.Image().Bounds()
	backdropBtn := ui.NewButton(ui.NewRectangle(m.app.Canvas, b.Dx(), b.Dy(), color.Transparent))
	m.backdropId = m.app.Touch.AddTapHandler(backdropBtn, func() {})
	_ = backdropBtn.AttachTo(m.backdrop)

	m.content.Mount()

	size := m.content.Size()
	m.panel.Add(ui.NewRectangle(m.app.Canvas, size.X, size.Y, color.Gray{Y: 200}))
}

func (m *Modal) Close() {
	m.backdrop.Clear()
	m.app.Touch.RemoveTapHandler(m.backdropId)
	m.panel.Clear()
	m.content.Unmount()
}

func (m *Modal) Render() []image.Rectangle {
	m.backdrop.Render(0, 0)
	dirtyRegs := m.panel.Render(500, 500)
	dirtyRegs = append(dirtyRegs, m.content.Render(500, 500)...)
	return ui.CombineRegions(dirtyRegs)
}
