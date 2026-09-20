package ui

type BaseGroup struct {
	Children []Renderer
}

func (bg *BaseGroup) Add(child Renderer) {
	bg.Children = append(bg.Children, child)
}
