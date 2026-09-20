package app

type Region struct {
	X0, Y0, X1, Y1 int
}

// TODO: rename to Displayer
type Renderer interface {
	Refresh([]Region) error
}

type Tappable interface {
	Contains(x, y int) bool
}

type TapDispatcher interface {
	AddTapHandler(btn Tappable, action func()) (id int)
	RemoveTapHandler(id int)
	Listen()
}
