//go:build js

package main

import (
	"image"
	"syscall/js"

	"github.com/gchumillas/prism/internal/app"
)

// IMPORTANT: these values must match those defined in docs/src/app.tsx
const (
	canvasWidth  = 468
	canvasHeight = 624
)

var (
	syncCanvas js.Func
	tapCanvas  js.Func
)

func platformInit() error {
	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	td := &TouchDispatcher{}
	ws := &Screen{canvas: canvas}
	app := &app.App{
		Touch:    td,
		Screen:   ws,
		Canvas:   canvas,
		FontSize: 10,
		DPI:      96,
	}

	syncCanvas = js.FuncOf(func(this js.Value, args []js.Value) any {
		js.CopyBytesToGo(canvas.Pix, args[0])
		ws.refreshCb = args[1]
		app.Run()
		return nil
	})

	tapCanvas = js.FuncOf(func(this js.Value, args []js.Value) any {
		x, y := args[0], args[1]
		td.tap(x.Int(), y.Int())
		return nil
	})

	js.Global().Set("prism", js.ValueOf(map[string]any{
		"syncCanvas": syncCanvas,
		"tapCanvas":  tapCanvas,
	}))

	select {}
}

// ---- Screen ----
type Screen struct {
	canvas    *image.RGBA
	refreshCb js.Value
}

func (s *Screen) Refresh(regs []app.Region) error {
	buf := js.Global().Get("Uint8Array").New(len(s.canvas.Pix))
	js.CopyBytesToJS(buf, s.canvas.Pix)
	for _, r := range regs {
		x0, y0 := js.ValueOf(r.X0), js.ValueOf(r.Y0)
		x1, y1 := js.ValueOf(r.X1), js.ValueOf(r.Y1)
		s.refreshCb.Invoke(buf, x0, y0, x1, y1)
	}
	return nil
}

// ---- TouchDispatcher ---- //
type Handler struct {
	id     int
	btn    app.Tappable
	action func()
}

type TouchDispatcher struct {
	handlers []Handler
	count    int
}

func (t *TouchDispatcher) tap(x, y int) {
	for _, h := range t.handlers {
		if h.btn.Contains(x, y) {
			h.action()
		}
	}
}

func (t *TouchDispatcher) AddTapHandler(btn app.Tappable, action func()) int {
	t.count++
	t.handlers = append(t.handlers, Handler{
		id:     t.count,
		btn:    btn,
		action: action,
	})
	return t.count
}

func (t *TouchDispatcher) RemoveTapHandler(id int) {
	for i, h := range t.handlers {
		if h.id == id {
			t.handlers = append(t.handlers[:i], t.handlers[i+1:]...)
			break
		}
	}
}
