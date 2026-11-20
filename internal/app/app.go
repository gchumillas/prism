package app

import "image"

// TODO: rename to Displayer
type Renderer interface {
	Refresh([]image.Rectangle) error
}

type Tappable interface {
	Contains(x, y int) bool
}

type TapDispatcher interface {
	AddTapHandler(btn Tappable, action func()) (id int)
	RemoveTapHandler(id int)
	Listen()
}
