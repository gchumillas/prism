package ui_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"

	_ "embed"
	_ "image/png"
)

type Spy[T ui.Widget] struct {
	Target   T
	rendered bool
}

func NewSpy[T ui.Widget](target T) *Spy[T] {
	return &Spy[T]{Target: target}
}

func (s *Spy[T]) TakeRendered() bool {
	r := s.rendered
	s.rendered = false
	return r
}

func (s *Spy[T]) Size() image.Point {
	return s.Target.Size()
}

func (s *Spy[T]) State() map[string]any {
	return s.Target.State()
}

func (s *Spy[T]) Render(x, y int) []image.Rectangle {
	s.rendered = true
	return s.Target.Render(x, y)
}

func TestVStack_Render(t *testing.T) {
	appWidth, appHeight := 320, 240

	t.Run("render all widgets", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)

		vs := ui.NewVStack(ctx)
		var shapes []*Spy[*ui.Rectangle]
		for i := range 3 {
			shape := NewSpy(ui.NewRectangle(ctx, 100, 50, color.Gray{Y: uint8(i)}))
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		// all shapes are rendered
		vs.Render(0, 0)
		for _, shape := range shapes {
			require.True(t, shape.TakeRendered())
		}

		// no shape is rendered
		vs.Render(0, 0)
		for _, shape := range shapes {
			require.False(t, shape.TakeRendered())
		}
	})

	t.Run("render only affected widget", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)

		vs := ui.NewVStack(ctx)
		var shapes []*Spy[*ui.LabelButton]
		for i := range 3 {
			shape := NewSpy(ui.NewLabelButton(ctx, fmt.Sprintf("Button %d", i)))
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		for _, shape := range shapes {
			shape.TakeRendered()
		}

		for i, shape := range shapes {
			shape.Target.Padding.X += 50
			vs.Render(0, 0)
			for j, sh := range shapes {
				require.Equal(t, j == i, sh.TakeRendered())
			}
		}
	})

	t.Run("render affected widget and below", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)

		vs := ui.NewVStack(ctx)
		var shapes []*Spy[*ui.LabelButton]
		for i := range 3 {
			shape := NewSpy(ui.NewLabelButton(ctx, fmt.Sprintf("Button %d", i)))
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		for _, shape := range shapes {
			shape.TakeRendered()
		}

		for i, shape := range shapes {
			shape.Target.Padding.Y += 50
			vs.Render(0, 0)
			for j, sh := range shapes {
				require.Equal(t, j >= i, sh.TakeRendered())
			}
		}
	})

	t.Run("render widgets below deleted", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)

		vs := ui.NewVStack(ctx)
		var shapes []*Spy[*ui.ImageButton]
		for i := range 3 {
			dat, err := os.ReadFile(fmt.Sprintf("testdata/icon%d.png", i+1))
			require.NoError(t, err)
			btn, err := ui.LoadImageButton(ctx, bytes.NewReader(dat))
			require.NoError(t, err)
			shape := NewSpy(btn)
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		for _, shape := range shapes {
			shape.TakeRendered()
		}

		for len(shapes) > 0 {
			vs.Remove(shapes[0])
			shapes = shapes[1:]
			vs.Render(0, 0)
			for _, sh := range shapes {
				require.True(t, sh.TakeRendered())
			}
		}
	})

	t.Run("render on change position", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)

		vs := ui.NewVStack(ctx)
		var shapes []*Spy[*ui.LabelButton]
		for i := range 3 {
			shape := NewSpy(ui.NewLabelButton(ctx, fmt.Sprintf("Button %d", i)))
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		for _, shape := range shapes {
			shape.TakeRendered()
		}

		// position changed and therefore all children must be rendered
		vs.Render(100, 0)
		for _, shape := range shapes {
			require.True(t, shape.TakeRendered())
		}

		// position remains and therefore all children must not be rendered
		vs.Render(100, 0)
		for _, shape := range shapes {
			require.False(t, shape.TakeRendered())
		}
	})

	t.Run("nested", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)

		vs1 := ui.NewVStack(ctx)
		vs2 := ui.NewVStack(ctx)
		vs1.Add(vs2)

		sh := NewSpy(ui.NewLabelButton(ctx, "Button 1"))
		vs2.Add(sh)

		vs1.Render(0, 0)
		require.True(t, sh.TakeRendered())

		sh.Target.Padding.X += 50
		vs1.Render(0, 0)
		require.True(t, sh.TakeRendered())
	})

	t.Run("nested change after steady state", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)

		vs1 := ui.NewVStack(ctx)
		vs2 := ui.NewVStack(ctx)
		vs1.Add(vs2)

		sh := NewSpy(ui.NewLabelButton(ctx, "Button 1"))
		vs2.Add(sh)

		vs1.Render(0, 0)
		require.True(t, sh.TakeRendered())

		// render until no widget changes (steady state)
		vs1.Render(0, 0)
		vs1.Render(0, 0)
		sh.TakeRendered()

		// a nested change must still be detected by the outer stack
		sh.Target.Padding.X += 50
		vs1.Render(0, 0)
		require.True(t, sh.TakeRendered())
	})

	t.Run("hit area follows widget after sibling removal", func(t *testing.T) {
		ctx := ui.NewContext(appWidth, appHeight)
		vs := ui.NewVStack(ctx)
		btn1 := ui.NewLabelButton(ctx, "Button")
		btn2 := ui.NewLabelButton(ctx, "Button")
		vs.Add(btn1)
		vs.Add(btn2)

		// btn2 is clickable at position p
		p := image.Pt(80, 115)
		vs.Render(0, 0)
		require.True(t, btn2.Contains(p.X, p.Y))

		// btn2 is no longer clickable at position p
		// because removing btn1 shifts btn2 upward
		vs.Remove(btn1)
		vs.Render(0, 0)
		require.False(t, btn2.Contains(p.X, p.Y))
	})
}
