//go:build !js

package platform

import (
	"errors"
	"image"
	"slices"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/screen"
	"github.com/gchumillas/prism/touch"
)

type Env struct {
	Touch    app.TapDispatcher
	Screen   app.Renderer
	Canvas   *image.RGBA
	FontSize float64
}

func LoadRuntime() (_ *Env, err error) {
	s, err := OpenScreen()
	if err != nil {
		return nil, err
	}

	td, err := OpenTouchDispatcher()
	if err != nil {
		return nil, err
	}

	return &Env{
		Touch:    td,
		Screen:   s,
		Canvas:   s.Canvas(),
		FontSize: 10,
	}, nil
}

// ---- Screen ----
type Screen struct {
	screen *screen.Screen
}

func OpenScreen() (s *Screen, err error) {
	s = &Screen{}
	if s.screen, err = screen.Open(); err != nil {
		return
	}
	return
}

func (s *Screen) Canvas() *image.RGBA {
	return s.screen.Canvas()
}

func (s *Screen) Refresh(regs []image.Rectangle) error {
	var rs []screen.Region
	for _, reg := range regs {
		rs = append(rs, screen.Region{
			X0: int(reg.Min.X), Y0: int(reg.Min.Y), X1: int(reg.Max.X), Y1: int(reg.Max.Y),
			Waveform: screen.WaveformGC16,
		})
	}
	return s.screen.RefreshRegions(rs)
}

// ---- TouchDispatcher ----
type Handler struct {
	id     int
	btn    app.Tappable
	action func()
}

type TouchDispatcher struct {
	touch    *touch.Touch
	handlers []Handler
	count    int
}

func OpenTouchDispatcher() (td *TouchDispatcher, err error) {
	td = &TouchDispatcher{}
	if td.touch, err = touch.Open(); err != nil {
		return
	}
	return
}

func (td *TouchDispatcher) AddTapHandler(btn app.Tappable, action func()) (id int) {
	td.count++
	td.handlers = append(td.handlers, Handler{
		id:     td.count,
		btn:    btn,
		action: action,
	})
	return td.count
}

func (td *TouchDispatcher) RemoveTapHandler(id int) error {
	for i, h := range td.handlers {
		if h.id == id {
			td.handlers = slices.Delete(td.handlers, i, i+1)
			return nil
		}
	}
	return errors.New("handler ID not found")
}

func (td *TouchDispatcher) Listen() {
	ts := &touch.TouchState{}
	var curPoint, prevPoint touch.TouchPoint
	for {
		if err := td.touch.Read(ts); err != nil {
			panic(err)
		}

		curPoint = ts.Points[0]
		if curPoint.ID < 0 && prevPoint.ID > 0 {
			handlers := slices.Clone(td.handlers)
			for i := len(handlers) - 1; i >= 0; i-- {
				handler := &handlers[i]
				if handler.btn.Contains(int(prevPoint.X), int(screen.Height-prevPoint.Y)) {
					handler.action()
					break
				}
			}
		}

		prevPoint = curPoint
	}
}
