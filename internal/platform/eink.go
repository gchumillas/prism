//go:build !js

package platform

import (
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

func (s *Screen) Refresh(regs []app.Region) error {
	var rs []screen.Region
	for _, reg := range regs {
		rs = append(rs, screen.Region{
			X0: int(reg.X0), Y0: int(reg.Y0), X1: int(reg.X1), Y1: int(reg.Y1),
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

func (td *TouchDispatcher) RemoveTapHandler(id int) {
	for i, h := range td.handlers {
		if h.id == id {
			td.handlers = slices.Delete(td.handlers, i, i+1)
			break
		}
	}
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
			for _, handler := range handlers {
				if handler.btn.Contains(int(prevPoint.X), int(screen.Height-prevPoint.Y)) {
					handler.action()
				}
			}
		}

		prevPoint = curPoint
	}
}
