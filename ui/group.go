package ui

type Group struct {
	Elems []Element
}

func (l *Group) Add(elem Element) {
	l.Elems = append(l.Elems, elem)
}
