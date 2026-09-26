package ui

import (
	"github.com/fogleman/gg"
)

const (
	defaultFontSize = 10
	defaultDPI      = 226
)

type Context struct {
	canvas   *gg.Context
	fontSize float64
}

type ContextOpt interface {
	applyContext(*Context)
}

// avail options
func (val fontSizeOpt) applyContext(x *Context) { x.fontSize = float64(val) }

// TODO: el canvas debería construirse internamente
func NewContext(canvas *gg.Context, opts ...ContextOpt) *Context {
	c := &Context{
		canvas:   canvas,
		fontSize: defaultFontSize,
	}
	for _, opt := range opts {
		opt.applyContext(c)
	}
	return c
}
