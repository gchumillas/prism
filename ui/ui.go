package ui

type Element interface {
	Size() (width, height int)
	Position() (x, y int)
	SetPosition(x, y int)
	Draw()
}
