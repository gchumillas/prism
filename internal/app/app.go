package app

type Region struct {
	X0, Y0, X1, Y1 float64
}

// TODO: rename to Displayer
type Renderer interface {
	Refresh([]Region) error
}

type Tappable interface {
	Contains(x, y float64) bool
}

type TapDispatcher interface {
	AddTapHandler(btn Tappable, action func()) (id int)
	RemoveTapHandler(id int)
	Listen()
}
