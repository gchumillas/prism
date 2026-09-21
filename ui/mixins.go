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

type BasePosition struct {
	x, y int
}

func (bp *BasePosition) Position() (x, y int) { return bp.x, bp.y }
func (bp *BasePosition) SetPosition(x, y int) { bp.x, bp.y = x, y }

type BaseSize struct {
	width, height int
}

func (bs *BaseSize) Size() (width, height int) { return bs.width, bs.height }
func (bs *BaseSize) SetSize(width, height int) { bs.width, bs.height = width, height }
