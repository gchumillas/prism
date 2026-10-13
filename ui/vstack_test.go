package ui_test

import (
	"bytes"
	"fmt"
	"hash/maphash"
	"image"
	"image/color"
	"image/draw"
	"os"
	"testing"

	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"

	_ "embed"
	_ "image/png"
)

type CustomWidget struct{ Label ui.Label }

func NewCustomWidget(c ui.Context, text string) *CustomWidget {
	return &CustomWidget{Label: *ui.NewLabel(c, text)}
}
func (cw *CustomWidget) Render(x, y int) []image.Rectangle { return cw.Label.Render(x, y) }
func (cw *CustomWidget) Size() image.Point                 { return cw.Label.Size() }
func (cw *CustomWidget) Footprint() []uint64               { return cw.Label.Footprint() }

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

func (s *Spy[T]) Size() image.Point   { return s.Target.Size() }
func (s *Spy[T]) Footprint() []uint64 { return s.Target.Footprint() }
func (s *Spy[T]) Render(x, y int) []image.Rectangle {
	s.rendered = true
	return s.Target.Render(x, y)
}

type Rect struct {
	canvas draw.Image
	image.Rectangle
	color color.Color
}

func NewRect(canvas draw.Image, x0, y0, x1, y1 int, c color.Color) *Rect {
	return &Rect{
		canvas:    canvas,
		Rectangle: image.Rect(x0, y0, x1, y1),
		color:     c,
	}
}

func (r *Rect) Size() image.Point {
	return r.Rectangle.Size()
}

func (r *Rect) Footprint() []uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, r.Dx())
	maphash.WriteComparable(&h, r.Dy())
	maphash.WriteComparable(&h, r.color)
	return []uint64{h.Sum64()}
}

func (r *Rect) Render(x, y int) []image.Rectangle {
	pos := image.Pt(x, y)
	size := r.Rectangle.Size()
	r.Rectangle = image.Rectangle{Min: pos, Max: pos.Add(size)}
	draw.Draw(r.canvas, r.Rectangle, image.NewUniform(r.color), image.Pt(0, 0), draw.Src)
	return []image.Rectangle{r.Rectangle}
}

type Rect1 struct {
	*Rect
}

func NewRect1(canvas draw.Image, x0, y0, x1, y1 int, c color.Color) *Rect1 {
	return &Rect1{Rect: NewRect(canvas, x0, y0, x1, y1, c)}
}

func TestVStack(t *testing.T) {
	// TODO: here we are testing VStack.Remove
	t.Run("can remove widgets", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)

		lbBtn := ui.NewLabelButton(env.ctx, "Click me!")
		vs.Add(lbBtn)
		require.Contains(t, vs.Children(), lbBtn)
		vs.Remove(lbBtn)
		require.NotContains(t, vs.Children(), lbBtn)
	})

	// TODO: we can create stackable custom widgets
	t.Run("can create stackable widgets", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		cw := NewCustomWidget(env.ctx, "custom ui")
		vs.Add(cw)
	})
}

func TestVStack_Render(t *testing.T) {
	intersectArea := func(r image.Rectangle, regs []image.Rectangle) int {
		area := 0
		for _, reg := range regs {
			s := reg.Intersect(r)
			area += s.Dx() * s.Dy()
		}
		return area
	}

	isRectInside := func(r image.Rectangle, regs []image.Rectangle) bool {
		return intersectArea(r, regs) == r.Dx()*r.Dy()
	}

	t.Run("render all widgets", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.ctx)
		var shapes []*Spy[*ui.Rectangle]
		for i := range 3 {
			shape := NewSpy(ui.NewRectangle(env.ctx, 100, 50, color.Gray{Y: uint8(i)}))
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
		env := NewEnv()

		vs := ui.NewVStack(env.ctx)
		var shapes []*Spy[*ui.Rectangle]
		for i := range 3 {
			shape := NewSpy(ui.NewRectangle(env.ctx, 100, 50, color.Gray{Y: uint8(i)}))
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		for _, shape := range shapes {
			shape.TakeRendered()
		}

		for i, shape := range shapes {
			shape.Target.Width += 50
			vs.Render(0, 0)
			for j, sh := range shapes {
				require.Equal(t, j == i, sh.TakeRendered())
			}
		}
	})

	t.Run("render affected widget and below", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.ctx)
		var shapes []*Spy[*ui.Rectangle]
		for i := range 3 {
			shape := NewSpy(ui.NewRectangle(env.ctx, 100, 50, color.Gray{Y: uint8(i)}))
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		for _, shape := range shapes {
			shape.TakeRendered()
		}

		for i, shape := range shapes {
			shape.Target.Height += 50
			vs.Render(0, 0)
			for j, sh := range shapes {
				require.Equal(t, j >= i, sh.TakeRendered())
			}
		}
	})

	t.Run("render widgets below deleted", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.ctx)
		var shapes []*Spy[*ui.ImageButton]
		for i := range 3 {
			dat, err := os.ReadFile(fmt.Sprintf("testdata/icon%d.png", i+1))
			require.NoError(t, err)
			btn, err := ui.LoadImageButton(env.ctx, bytes.NewReader(dat))
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
		env := NewEnv()

		vs := ui.NewVStack(env.ctx)
		var shapes []*Spy[*ui.LabelButton]
		for i := range 3 {
			shape := NewSpy(ui.NewLabelButton(env.ctx, fmt.Sprintf("Button %d", i)))
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
		env := NewEnv()

		vs1 := ui.NewVStack(env.ctx)
		vs2 := ui.NewVStack(env.ctx)
		vs1.Add(vs2)

		sh := NewSpy(ui.NewRectangle(env.ctx, 100, 50, color.White))
		vs2.Add(sh)

		vs1.Render(0, 0)
		require.True(t, sh.TakeRendered())

		sh.Target.Width += 50
		vs1.Render(0, 0)
		require.True(t, sh.TakeRendered())
	})

	t.Run("nested change after steady state", func(t *testing.T) {
		env := NewEnv()

		vs1 := ui.NewVStack(env.ctx)
		vs2 := ui.NewVStack(env.ctx)
		vs1.Add(vs2)

		sh := NewSpy(ui.NewRectangle(env.ctx, 100, 50, color.White))
		vs2.Add(sh)

		vs1.Render(0, 0)
		require.True(t, sh.TakeRendered())

		// render until no widget changes (steady state)
		vs1.Render(0, 0)
		vs1.Render(0, 0)
		sh.TakeRendered()

		// a nested change must still be detected by the outer stack
		sh.Target.Width += 50
		vs1.Render(0, 0)
		require.True(t, sh.TakeRendered())
	})

	t.Run("hit area follows widget after sibling removal", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		btn1 := ui.NewLabelButton(env.ctx, "Button")
		btn2 := ui.NewLabelButton(env.ctx, "Button")
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

	t.Run("remove from list of different widgets", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		var shapes []*ui.Rectangle
		for i := range uint8(2) {
			r := ui.NewRectangle(env.ctx, 100, 50, color.Gray{Y: i})
			shapes = append(shapes, r)
			vs.Add(r)
		}
		vs.Render(0, 0)

		sh0 := shapes[0]
		vs.Remove(sh0)
		regs := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 100)}, regs)
	})

	t.Run("remove from list of identical widgets", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		var shapes []*ui.Rectangle
		for range uint8(3) {
			r := ui.NewRectangle(env.ctx, 100, 50, color.White)
			shapes = append(shapes, r)
			vs.Add(r)
		}
		vs.Render(0, 0)

		sh0 := shapes[0]
		vs.Remove(sh0)
		regs := vs.Render(0, 0)
		// only returns the region of the last rectangle
		// since the first rectangle is identical to the second and the third,
		// there is no need to repaint those regions
		require.Equal(t, []image.Rectangle{image.Rect(0, 100, 100, 150)}, regs)
	})

	t.Run("render different widget sizes", func(t *testing.T) {
		env := NewEnv()
		recs := []*Rect{
			NewRect(env.canvas, 0, 0, 100, 50, color.Gray{Y: 25}),  // landscape
			NewRect(env.canvas, 0, 50, 50, 100, color.Gray{Y: 50}), // square
		}
		vs := ui.NewVStack(env.ctx)
		for _, rec := range recs {
			vs.Add(rec)
		}
		vs.Render(0, 0)

		r0 := recs[0].Rectangle
		r1 := recs[1].Rectangle
		vs.Remove(recs[0])
		regs := vs.Render(0, 0)
		require.True(t, isRectInside(r0, regs))
		require.True(t, isRectInside(recs[1].Rectangle, regs))
		require.True(t, isRectInside(r1, regs))
	})

	t.Run("different widget with the same footprint", func(t *testing.T) {
		env := NewEnv()
		recs := []ui.Widget{
			NewRect(env.canvas, 0, 0, 100, 50, color.White),
			NewRect1(env.canvas, 0, 50, 100, 100, color.White),
		}
		vs := ui.NewVStack(env.ctx, ui.WithGap(10))
		for _, rec := range recs {
			vs.Add(rec)
		}
		regs := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 50),
			image.Rect(0, 60, 100, 110),
		}, regs)

		vs.Remove(recs[0])
		regs = vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 50),   // must be repainted, since Rect1 and Rect are of different types
			image.Rect(0, 60, 100, 110), // and the abandoned region must also be repainted
		}, regs)
	})

	t.Run("nested vstack child replaced by different type", func(t *testing.T) {
		env := NewEnv()

		vs1 := ui.NewVStack(env.ctx)
		vs2 := NewSpy(ui.NewVStack(env.ctx))
		vs1.Add(vs2)

		sh := NewRect(env.canvas, 0, 0, 100, 50, color.Black)
		vs2.Target.Add(sh)
		vs1.Render(0, 0)
		require.True(t, vs2.TakeRendered())

		// replace child with another one with identical footprint and different type
		vs2.Target.Remove(sh)
		vs2.Target.Add(NewRect1(env.canvas, 0, 0, 100, 50, color.Black))
		vs1.Render(0, 0)
		require.True(t, vs2.TakeRendered())
	})
}
