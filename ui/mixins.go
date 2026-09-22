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

func (bp *BasePosition) X() int       { return bp.x }
func (bp *BasePosition) Y() int       { return bp.y }
func (bp *BasePosition) SetX(val int) { bp.x = val }
func (bp *BasePosition) SetY(val int) { bp.y = val }

type BaseSize struct {
	width, height int
}

func (bs *BaseSize) Width() int        { return bs.width }
func (bs *BaseSize) Height() int       { return bs.height }
func (bs *BaseSize) SetWidth(val int)  { bs.width = val }
func (bs *BaseSize) SetHeight(val int) { bs.height = val }
