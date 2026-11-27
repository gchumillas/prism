//go:build regression

package ui_test

import (
	"hash/maphash"
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/fogleman/gg"
	"github.com/gchumillas/prism/internal/test"
	"github.com/gchumillas/prism/ui"
	"github.com/stretchr/testify/require"
)

type Square struct {
	ctx   *gg.Context
	width int
	bSize int
}

func NewSquare(canvas *image.RGBA, width, bSize int) *Square {
	ctx := gg.NewContextForRGBA(canvas)
	return &Square{ctx: ctx, width: width, bSize: bSize}
}

func (sq *Square) Footprint() []uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	maphash.WriteComparable(&h, sq.width)
	maphash.WriteComparable(&h, sq.bSize)
	return []uint64{h.Sum64()}
}

func (sq *Square) Render(x, y int) []image.Rectangle {
	sq.ctx.SetColor(color.White)
	sq.ctx.SetLineWidth(float64(sq.bSize))
	sq.ctx.DrawLine(float64(x), float64(y), float64(x+sq.width), float64(y))
	sq.ctx.DrawLine(float64(x), float64(y+sq.width), float64(x+sq.width), float64(y+sq.width))
	sq.ctx.DrawLine(float64(x), float64(y), float64(x), float64(y+sq.width))
	sq.ctx.DrawLine(float64(x+sq.width), float64(y), float64(x+sq.width), float64(y+sq.width))
	sq.ctx.Stroke()

	safeMargin := sq.bSize/2 + 1
	dirtyReg := image.Rectangle{
		Min: image.Pt(x-safeMargin, y-safeMargin),
		Max: image.Pt(x+sq.width+safeMargin, y+sq.width+safeMargin),
	}

	return []image.Rectangle{dirtyReg}
}

func (sq *Square) Size() image.Point {
	return image.Pt(sq.width, sq.width)
}

// https://github.com/gchumillas/prism/issues/110
func TestIssue110(t *testing.T) {
	width := 50
	x0, y0 := 100, 50
	x1, y1 := x0+50, y0+50
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// internal stack
	hs1 := ui.NewHStack()
	r := NewRectangle(c, width, width, color.White)
	hs1.Add(r)

	// render a external stack that contains a internal stack
	hs := ui.NewHStack()
	_ = hs1.AttachTo(hs)
	hs.Render(x0, y0)

	// remove the internal stack, render and clear "dirty regions"
	_ = hs1.DetachFrom(hs)
	dirtyRegs := hs.Render(x0, y0)
	for _, reg := range dirtyRegs {
		draw.Draw(c, reg, image.NewUniform(color.Transparent), image.Point{}, draw.Src)
	}
	test.RectIsColor(t, c, x0, y0, x1, y1, color.Transparent) // has been cleared correctly

	// add the internal stack again and render
	_ = hs1.AttachTo(hs)
	hs.Render(x0, y0)
	test.RectIsColor(t, c, x0+10, y0+10, x1-10, y1-10, color.White) // has been rendered correctly

	// remove the internal stack again, repaint and clear "dirty regions"
	hs1.DetachFrom(hs)
	dirtyRegs = hs.Render(x0, y0)
	for _, reg := range dirtyRegs {
		draw.Draw(c, reg, image.NewUniform(color.Transparent), image.Point{}, draw.Src)
	}
	test.RectIsColor(t, c, x0, y0, x1, y1, color.Transparent) // has been cleared correctly
}

// https://github.com/gchumillas/prism/issues/73
func TestIssue73(t *testing.T) {
	t.Run("dirty regions cover all pixels", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		sq1 := NewSquare(c, 50, 10) // square with a border
		sq2 := NewSquare(c, 50, 0)  // square without a border
		hs := ui.NewHStack(ui.WithGap(50))
		hs.Add(sq1)
		hs.Add(sq2)
		hs.Render(10, 10)

		// clear the squares and repaint
		hs.Remove(sq1)
		hs.Remove(sq2)
		dirtyRegs := hs.Render(10, 10)

		// clear "dirty regions"
		for _, reg := range dirtyRegs {
			draw.Draw(c, reg, image.NewUniform(color.Transparent), image.Pt(0, 0), draw.Src)
		}

		// the canvas should be completely clean
		b := c.Bounds()
		test.RectIsColor(t, c, b.Min.X, b.Min.Y, b.Max.X, b.Max.Y, color.Transparent)
	})

	t.Run("prev regions are preserved", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		hs := ui.NewHStack()

		square := NewRectangle(c, 20, 20, color.White)
		hs.Add(square)
		dirtyRegs := hs.Render(0, 0)
		// "dirtyRegs" contains the square and only it
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 20, 20),
		}, dirtyRegs)

		hs.Remove(square)
		landscape := NewRectangle(c, 30, 10, color.White)
		hs.Add(landscape)
		dirtyRegs = hs.Render(0, 0)
		// "dirtyRegs" contains the square and the landscape and only them
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 30, 10),
			image.Rect(0, 10, 20, 20),
		}, dirtyRegs)

		hs.Remove(landscape)
		portrait := NewRectangle(c, 10, 30, color.Black)
		hs.Add(portrait)
		// "dirtyRegs" contains the landscape and the portrait and only them
		// it doesn't contain the square
		dirtyRegs = hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 30, 10),
			image.Rect(0, 10, 10, 30),
		}, dirtyRegs)
	})
}

// https://github.com/gchumillas/prism/issues/113
func TestIssue113(t *testing.T) {
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))
	hs1 := ui.NewHStack()
	hs2 := ui.NewHStack()
	a := NewRectangle(c, 50, 50, color.Gray{Y: 75})
	b := NewRectangle(c, 50, 50, color.Gray{Y: 150})
	hs2.Add(a)
	hs2.Add(b)
	hs2.AttachTo(hs1)

	// An inner HStack contains rectangles A and B.
	// The first render returns the regions of A and B, which are cached.
	regs := hs1.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(0, 0, 100, 50),
	}, regs)

	// B is removed and rendered again. The inner HStack returns only B's old
	// region, because that's the only thing that changed. Only that is cached.
	hs2.Remove(b)
	regs = hs1.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(50, 0, 100, 50),
	}, regs)

	// The inner HStack is detached. prev.regs is cleared, and contains only B.
	hs2.DetachFrom(hs1)
	regs = hs1.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(0, 0, 50, 50),
	}, regs)
}

// https://github.com/gchumillas/prism/issues/114
func TestIssue114(t *testing.T) {
	t.Run("partially change nested stack", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		hs1 := ui.NewHStack()
		hs2 := ui.NewHStack()
		var recs []*Rectangle
		for i := range uint8(3) {
			r := NewRectangle(c, 50, 50, color.Gray{Y: 50 + 50*i})
			recs = append(recs, r)
			hs2.Add(r)
		}
		hs2.AttachTo(hs1)
		hs1.Render(0, 0)

		recs[1].height = 100
		regs := hs1.Render(0, 0)
		// only recs[1] needs to be refreshed
		// however regs is [(0,0)-(150,50) (50,50)-(100,100)]
		require.Equal(t, []image.Rectangle{
			image.Rect(50, 0, 100, 100),
		}, regs)
	})

	// https://github.com/gchumillas/prism/pull/119#discussion_r4364201290
	t.Run("replace widget with stack", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		hs1 := ui.NewHStack()
		big := NewRectangle(c, 50, 50, color.White)
		hs1.Add(big)
		hs1.Render(0, 0)

		// replace the rectangle with a nested stack that has a smaller content
		hs1.Remove(big)
		hs2 := ui.NewHStack()
		hs2.Add(NewRectangle(c, 20, 20, color.White))
		hs2.AttachTo(hs1)
		regs := hs1.Render(0, 0)

		// the old 50x50 region must be refreshed, otherwise stale pixels remain
		for _, reg := range regs {
			draw.Draw(c, reg, image.NewUniform(color.Transparent), image.Point{}, draw.Src)
		}
		test.RectIsColor(t, c, 0, 0, 50, 50, color.Transparent)
	})
}

// https://github.com/gchumillas/prism/pull/120#discussion_r4368093650
func TestPull120(t *testing.T) {
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))
	hs := ui.NewHStack()
	hs1 := ui.NewHStack()
	hs1.AttachTo(hs)

	regs := hs.Render(0, 0)
	require.Empty(t, regs)

	r := NewRectangle(c, 50, 50, color.White)
	hs1.Add(r)
	regs = hs.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(0, 0, 50, 50),
	}, regs)

	hs1.DetachFrom(hs)
	regs = hs.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(0, 0, 50, 50),
	}, regs)
}

// https://github.com/gchumillas/prism/issues/124
func TestIssue124(t *testing.T) {
	t.Run("vertical", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		vs := ui.NewVStack()

		vs1 := ui.NewVStack()
		vs1.Add(NewRectangle(c, 100, 100, color.White))

		vs2 := ui.NewVStack()
		vs2.Add(NewRectangle(c, 20, 20, color.White))

		vs1.AttachTo(vs)
		regs := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 100),
		}, regs)

		vs1.DetachFrom(vs)
		err := vs2.AttachTo(vs)
		require.ErrorIs(t, err, ui.ErrNeedsRender)
	})
	t.Run("horizontal", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		hs := ui.NewHStack()

		hs1 := ui.NewHStack()
		hs1.Add(NewRectangle(c, 100, 100, color.White))

		hs2 := ui.NewHStack()
		hs2.Add(NewRectangle(c, 20, 20, color.White))

		hs1.AttachTo(hs)
		regs := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 100),
		}, regs)

		hs1.DetachFrom(hs)
		err := hs2.AttachTo(hs)
		require.ErrorIs(t, err, ui.ErrNeedsRender)
	})
	t.Run("detached button", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		vs := ui.NewVStack()

		btn := ui.NewButton(NewRectangle(c, 100, 100, color.White))
		vs2 := ui.NewVStack()
		vs2.Add(NewRectangle(c, 20, 20, color.White))

		btn.AttachTo(vs)
		regs := vs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 100),
		}, regs)

		btn.DetachFrom(vs)
		err := vs2.AttachTo(vs)
		require.ErrorIs(t, err, ui.ErrNeedsRender)
	})
}

// https://github.com/gchumillas/prism/issues/126
func TestIssue126(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		hs := ui.NewHStack()

		r := NewRectangle(c, 50, 50, color.White)

		hs1 := ui.NewHStack()
		hs1.Add(r)
		hs1.AttachTo(hs)

		hs2 := ui.NewHStack()
		hs2.Add(r)
		hs2.AttachTo(hs)

		hs.Render(0, 0)
		hs1.DetachFrom(hs)
		regs := hs.Render(0, 0)
		// only the missing square needs to be refreshed
		require.Equal(t, []image.Rectangle{
			image.Rect(50, 0, 100, 50),
		}, regs)
	})

	// https://github.com/gchumillas/prism/pull/127#discussion_r4372638791
	t.Run("case 2", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		stacks := []*ui.Stack{ui.NewHStack()}
		for _, st := range stacks {
			r := NewRectangle(c, 50, 50, color.White)
			st.Add(r)

			regs := st.Render(0, 0)
			require.Equal(t, []image.Rectangle{
				image.Rect(0, 0, 50, 50),
			}, regs)

			st.Remove(r)
			regs = st.Render(100, 100)
			require.Equal(t, []image.Rectangle{
				image.Rect(0, 0, 50, 50),
			}, regs)
		}
	})

	// https://github.com/gchumillas/prism/pull/127#discussion_r4372638811
	t.Run("case 3", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		st := ui.NewHStack()

		st1 := ui.NewHStack()
		red := NewRectangle(c, 50, 50, color.RGBA{R: 255})
		st1.Add(red)
		st1.AttachTo(st)

		st2 := ui.NewHStack()
		blue := NewRectangle(c, 50, 50, color.RGBA{B: 255})
		st2.Add(blue)
		st2.AttachTo(st)

		regs := st.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 50),
		}, regs)

		st1.DetachFrom(st)
		regs = st.Render(0, 0)
		// should return the joint shadow of st1 and st2,
		// since they contain squares of different colors
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 100, 50),
		}, regs)
	})

	// https://github.com/gchumillas/prism/pull/127#discussion_r4372638832
	t.Run("case 4", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		st := ui.NewHStack()
		r := NewRectangle(c, 50, 50, color.Black)
		st.Add(r)
		st1 := ui.NewHStack()
		st1.Add(NewRectangle(c, 50, 50, color.White))
		st1.AttachTo(st)
		st.Render(0, 0)
		r.width = 100
		regs := st.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 150, 50),
		}, regs)
	})

	// https://github.com/gchumillas/prism/pull/129#discussion_r4579982636
	t.Run("case 5", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		hs := ui.NewHStack()
		vs := ui.NewVStack()
		vs.Add(NewRectangle(c, 50, 50, color.White))
		vs.Add(NewRectangle(c, 50, 50, color.Black))
		_ = vs.AttachTo(hs)
		regs := hs.Render(0, 0)
		require.Equal(t, []image.Rectangle{
			image.Rect(0, 0, 50, 100),
		}, regs)
	})

	// https://github.com/gchumillas/prism/pull/129#discussion_r4579982654
	t.Run("case 6", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		hs := ui.NewHStack()
		hs1 := ui.NewHStack()
		hs1.PaddingY = 5
		hs1.Add(NewRectangle(c, 50, 50, color.Black))
		_ = hs1.AttachTo(hs)
		hs2 := ui.NewHStack()
		hs2.PaddingY = 5
		hs2.Add(NewRectangle(c, 50, 50, color.Black))
		_ = hs2.AttachTo(hs)
		size := hs.Size()
		require.Equal(t, image.Pt(100, 60), size)

		vs := ui.NewVStack()
		vs1 := ui.NewVStack()
		vs1.PaddingX = 5
		vs1.Add(NewRectangle(c, 50, 50, color.Black))
		_ = vs1.AttachTo(vs)
		vs2 := ui.NewVStack()
		vs2.PaddingX = 5
		vs2.Add(NewRectangle(c, 50, 50, color.Black))
		_ = vs2.AttachTo(vs)
		size1 := vs.Size()
		require.Equal(t, image.Pt(60, 100), size1)
	})

	// https://github.com/gchumillas/prism/pull/129#discussion_r4598996061
	t.Run("case 7", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		hs := ui.NewHStack()
		hs1 := ui.NewHStack()
		hs1.Add(NewRectangle(c, 50, 50, color.White))
		hs1.AttachTo(hs)
		vs1 := ui.NewVStack()
		vs1.Add(NewRectangle(c, 50, 50, color.White))
		vs1.AttachTo(hs)
		size := hs.Size()
		require.Equal(t, image.Pt(100, 50), size)

		vs := ui.NewVStack()
		vs2 := ui.NewVStack()
		vs2.Add(NewRectangle(c, 50, 50, color.White))
		vs2.AttachTo(vs)
		hs2 := ui.NewHStack()
		hs2.Add(NewRectangle(c, 50, 50, color.White))
		hs2.AttachTo(vs)
		size1 := vs.Size()
		require.Equal(t, image.Pt(50, 100), size1)
	})
}

// https://github.com/gchumillas/prism/pull/132
func TestPull132(t *testing.T) {
	// https://github.com/gchumillas/prism/pull/132#discussion_r4621777577
	t.Run("case 1", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		hs := ui.NewHStack()
		hs.PaddingX = 5
		hs.Add(NewRectangle(c, 50, 50, color.Black))
		size := hs.Size()
		require.Equal(t, image.Pt(60, 50), size)
	})

	// https://github.com/gchumillas/prism/pull/132#discussion_r4621777666
	t.Run("case 2", func(t *testing.T) {
		c := image.NewRGBA(image.Rect(0, 0, 320, 240))
		vs := ui.NewVStack()
		vs.PaddingY = 5
		vs.Add(NewRectangle(c, 50, 50, color.Black))
		size1 := vs.Size()
		require.Equal(t, image.Pt(50, 60), size1)
	})
}

// https://github.com/gchumillas/prism/issues/130
func TestIssue130(t *testing.T) {
	c := image.NewRGBA(image.Rect(0, 0, 320, 240))
	hs := ui.NewHStack()
	hs1 := ui.NewHStack()
	previous := NewRectangle(c, 50, 50, color.Gray{Y: 100})
	removed := NewRectangle(c, 50, 50, color.Black)
	hs.Add(previous)
	hs1.Add(removed)
	hs1.AttachTo(hs)
	hs.Render(0, 0)

	hs1.Remove(removed)
	regs := hs1.Render(50, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(50, 0, 100, 50),
	}, regs)

	regs = hs.Render(0, 0)
	require.Equal(t, []image.Rectangle{
		image.Rect(50, 0, 100, 50),
	}, regs)
}
