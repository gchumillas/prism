package ui_test

import (
	"image"
	"image/color"
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
}

func TestVStack_Add(t *testing.T) {
	t.Run("custom widget", func(t *testing.T) {
		env := NewEnv()
		vs := ui.NewVStack(env.ctx)
		cw := NewCustomWidget(env.ctx, "custom ui")
		vs.Add(cw)
		require.Contains(t, vs.Children(), cw)
	})
}

func TestVStack_Remove(t *testing.T) {
	env := NewEnv()
	vs := ui.NewVStack(env.ctx)
	r := NewRect(env.canvas, 0, 0, 100, 50, color.White)
	vs.Add(r)
	require.Contains(t, vs.Children(), r)
	vs.Remove(r)
	require.NotContains(t, vs.Children(), r)
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
		for i := range 3 {
			shape := NewRect(env.canvas, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			vs.Add(shape)
		}

		// all shapes are rendered
		regs := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 150)}, regs)

		// no shape is rendered
		regs = vs.Render(0, 0)
		require.Empty(t, regs)
	})

	t.Run("render only affected widget", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.ctx)
		var shapes []*Rect
		for i := range 3 {
			shape := NewRect(env.canvas, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		shapes[1].Max.X += 50
		regs := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 50, 150, 100)}, regs)
	})

	t.Run("render affected widget and below", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.ctx)
		var shapes []*Rect
		for i := range 3 {
			shape := NewRect(env.canvas, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		shapes[1].Max.Y += 50
		regs := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 50, 100, 200)}, regs)
	})

	t.Run("render widgets below deleted", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.ctx)
		var shapes []*Rect
		for i := range 3 {
			shape := NewRect(env.canvas, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			shapes = append(shapes, shape)
			vs.Add(shape)
		}

		vs.Render(0, 0)
		vs.Remove(shapes[0])
		regs := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 150)}, regs)
	})

	t.Run("render on change position", func(t *testing.T) {
		env := NewEnv()

		vs := ui.NewVStack(env.ctx)
		for i := range 3 {
			shape := NewRect(env.canvas, 0, 0, 100, 50, color.Gray{Y: uint8(i)})
			vs.Add(shape)
		}

		vs.Render(0, 0)
		regs := vs.Render(100, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 200, 150)}, regs)
	})

	t.Run("nested", func(t *testing.T) {
		env := NewEnv()

		vs1 := ui.NewVStack(env.ctx)
		vs2 := ui.NewVStack(env.ctx)
		vs1.Add(vs2)

		sh := NewRect(env.canvas, 0, 0, 100, 50, color.White)
		vs2.Add(sh)

		regs := vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 50)}, regs)

		sh.Max.X += 50
		regs = vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 150, 50)}, regs)
	})

	t.Run("nested change after steady state", func(t *testing.T) {
		env := NewEnv()

		vs1 := ui.NewVStack(env.ctx)
		vs2 := ui.NewVStack(env.ctx)
		vs1.Add(vs2)

		sh := NewRect(env.canvas, 0, 0, 100, 50, color.White)
		vs2.Add(sh)

		regs := vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 50)}, regs)

		// render until no widget changes (steady state)
		vs1.Render(0, 0)
		vs1.Render(0, 0)

		// a nested change must still be detected by the outer stack
		sh.Max.X += 50
		regs = vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 150, 50)}, regs)
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
		var shapes []*Rect
		for i := range uint8(2) {
			r := NewRect(env.canvas, 0, 0, 100, 50, color.Gray{Y: i})
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
		var shapes []*Rect
		for range uint8(3) {
			r := NewRect(env.canvas, 0, 0, 100, 50, color.White)
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
		vs2 := ui.NewVStack(env.ctx)
		vs1.Add(vs2)

		sh := NewRect(env.canvas, 0, 0, 100, 50, color.Black)
		vs2.Add(sh)
		regs := vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 50)}, regs)

		// replace child with another one with identical footprint and different type
		vs2.Remove(sh)
		vs2.Add(NewRect1(env.canvas, 0, 0, 100, 50, color.Black))
		regs = vs1.Render(0, 0)
		require.Equal(t, []image.Rectangle{image.Rect(0, 0, 100, 50)}, regs)
	})
}
