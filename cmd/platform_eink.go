//go:build !js

package main

import (
	"image/draw"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/screen"
	"github.com/gchumillas/prism/touch"
)

func platformInit() error {
	s, err := openScreen()
	if err != nil {
		return err
	}
	defer s.close()

	td, err := openTouchDispatcher()
	if err != nil {
		return err
	}

	app := &app.App{
		Touch:    td,
		Screen:   s,
		Canvas:   s.canvas(),
		FontSize: 10,
		DPI:      226,
	}
	if err := app.Run(); err != nil {
		panic(err)
	}

	td.listenEvents()
	return nil
}

// ---- eInkScreen ----
type einkScreen struct {
	screen *screen.Screen
}

func openScreen() (s *einkScreen, err error) {
	s = &einkScreen{}
	if s.screen, err = screen.Open(); err != nil {
		return
	}
	return
}

func (s *einkScreen) close() {
	s.screen.Close()
}

func (s *einkScreen) canvas() draw.Image {
	return s.screen.Canvas()
}

func (s *einkScreen) Refresh(regs []app.Region) error {
	var rs []screen.Region
	for _, reg := range regs {
		rs = append(rs, screen.Region{
			X0: reg.X0, Y0: reg.Y0, X1: reg.X1, Y1: reg.Y1,
			Waveform: screen.WaveformGC16,
		})
	}
	return s.screen.RefreshRegions(rs)
}

// ---- touchDispatcher ----
type handler struct {
	id     int
	btn    app.Tappable
	action func()
}

type touchDispatcher struct {
	touch    *touch.Touch
	handlers []handler
	count    int
}

func openTouchDispatcher() (td *touchDispatcher, err error) {
	td = &touchDispatcher{}
	if td.touch, err = touch.Open(); err != nil {
		return
	}
	return
}

func (td *touchDispatcher) AddTapHandler(btn app.Tappable, action func()) (id int) {
	td.count++
	td.handlers = append(td.handlers, handler{id: td.count, btn: btn, action: action})
	return td.count
}

func (td *touchDispatcher) RemoveTapHandler(id int) {
	for i, h := range td.handlers {
		if h.id == id {
			td.handlers = append(td.handlers[:i], td.handlers[i+1:]...)
			break
		}
	}
}

func (td *touchDispatcher) listenEvents() {
	ts := &touch.TouchState{}
	var curPoint, prevPoint touch.TouchPoint
	for {
		if err := td.touch.Read(ts); err != nil {
			panic(err)
		}

		curPoint = ts.Points[0]
		if curPoint.ID < 0 && prevPoint.ID > 0 {
			for _, handler := range td.handlers {
				x, y := int(prevPoint.X), int(screen.Height-prevPoint.Y)
				if handler.btn.Contains(x, y) {
					handler.action()
				}
			}
		}

		prevPoint = curPoint
	}
}
