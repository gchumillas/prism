package ui

type BaseGroup struct {
	Children []Renderer
}

func (bg *BaseGroup) Add(child Renderer) {
	bg.Children = append(bg.Children, child)
}

type BasePosition struct {
	X, Y int
}

func (bp *BasePosition) Position() (x, y int) { return bp.X, bp.Y }
func (bp *BasePosition) SetPosition(x, y int) { bp.X, bp.Y = x, y }
