package ui

type BaseGroup struct {
	// TODO: why Children is public?
	Children []Renderer
}

func (bg *BaseGroup) Add(child Renderer) {
	bg.Children = append(bg.Children, child)
}

func (b *BaseGroup) Remove(child Renderer) {
	for i, c := range b.Children {
		if c == child {
			b.Children = append(b.Children[:i], b.Children[i+1:]...)
			break
		}
	}
}

type BasePosition struct {
	X, Y int
}

func (bp *BasePosition) Position() (x, y int) { return bp.X, bp.Y }
func (bp *BasePosition) SetPosition(x, y int) { bp.X, bp.Y = x, y }

type BaseSize struct {
	Width, Height int
}

func (bs *BaseSize) Size() (width, height int) { return bs.Width, bs.Height }
func (bs *BaseSize) SetSize(width, height int) { bs.Width, bs.Height = width, height }
