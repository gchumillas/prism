package ui

import (
	"bytes"
	"encoding/json"
	"image"
	"slices"
)

func MergeRegions(regs []image.Rectangle) []image.Rectangle {
	type segment = [2]int

	var xcoords, ycoords []int
	for _, reg := range regs {
		xcoords = append(xcoords, reg.Min.X, reg.Max.X)
		ycoords = append(ycoords, reg.Min.Y, reg.Max.Y)
	}
	slices.Sort(xcoords)
	xcoords = slices.Compact(xcoords)
	slices.Sort(ycoords)
	ycoords = slices.Compact(ycoords)

	// horizontal segments
	var hSegments []segment
	for i := 1; i < len(xcoords); i++ {
		hSegments = append(hSegments, [2]int{xcoords[i-1], xcoords[i]})
	}

	// vertical segments
	var vSegments []segment
	for i := 1; i < len(ycoords); i++ {
		vSegments = append(vSegments, [2]int{ycoords[i-1], ycoords[i]})
	}

	table := make([][]*image.Rectangle, len(vSegments))
	for i, vs := range vSegments {
		row := make([]*image.Rectangle, len(hSegments))
		for j, hs := range hSegments {
			r := image.Rect(hs[0], vs[0], hs[1], vs[1])
			if i := slices.IndexFunc(regs, func(s image.Rectangle) bool {
				return r.In(s)
			}); i < 0 {
				continue
			}
			row[j] = &r
		}
		table[i] = row
	}

	// join adjacent horizontal cells
	for i := range vSegments {
		for j0, j := 0, 1; j < len(hSegments); j++ {
			cell0, cell := table[i][j0], table[i][j]
			if cell0 != nil && cell != nil && cell0.Dy() == cell.Dy() {
				cell0.Max = cell.Max
				table[i][j] = nil
				continue
			}
			j0 = j
		}
	}

	// join adjacent vertical cells
	for i := range hSegments {
		for j0, j := 0, 1; j < len(vSegments); j++ {
			cell0, cell := table[j0][i], table[j][i]
			if cell0 != nil && cell != nil && cell0.Dx() == cell.Dx() {
				cell0.Max = cell.Max
				table[j][i] = nil
				continue
			}
			j0 = j
		}
	}

	var result []image.Rectangle
	for _, row := range table {
		for _, cell := range row {
			if cell != nil {
				result = append(result, *cell)
			}
		}
	}
	return result
}

// used by "Organizers", like VStack, to determine if a node must be rendered
func serializeNode(x, y int, node Renderer) []byte {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(struct {
		X, Y int
		N    Renderer
	}{X: x, Y: y, N: node}); err != nil {
		panic(err)
	}
	return b.Bytes()
}
