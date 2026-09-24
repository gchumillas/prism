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
	position Point
}

func (bp *BasePosition) X() float64       { return bp.position.X }
func (bp *BasePosition) Y() float64       { return bp.position.Y }
func (bp *BasePosition) SetX(val float64) { bp.position.X = val }
func (bp *BasePosition) SetY(val float64) { bp.position.Y = val }

type BaseSize struct {
	size Point
}

func (bs *BaseSize) Width() float64        { return bs.size.X }
func (bs *BaseSize) Height() float64       { return bs.size.Y }
func (bs *BaseSize) SetWidth(val float64)  { bs.size.X = val }
func (bs *BaseSize) SetHeight(val float64) { bs.size.Y = val }
