package ui

type BaseGroup struct {
	children []Renderer
}

func (bg *BaseGroup) Children() []Renderer {
	return bg.children
}

func (bg *BaseGroup) Add(child Renderer) {
	bg.children = append(bg.children, child)
}

func (b *BaseGroup) Remove(child Renderer) {
	for i, c := range b.children {
		if c == child {
			b.children = append(b.children[:i], b.children[i+1:]...)
			break
		}
	}
}
