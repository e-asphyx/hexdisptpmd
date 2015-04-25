package main

const (
	Sqrt3 = 1.73205080756887729352744634150587236694280525381038062805580698 //A002194
	Size  = 1.0
)

type HexLayout struct {
	EdgeLen       int
	InvertOddRows bool
}

func (layout *HexLayout) NumCells() int {
	return 3*layout.EdgeLen*layout.EdgeLen - 3*layout.EdgeLen + 1
}

func (layout *HexLayout) rowLen(row int) int {
	if row < layout.EdgeLen {
		return layout.EdgeLen + row
	} else {
		return layout.EdgeLen*3 - row - 2
	}
}

func (layout *HexLayout) cellPos(n int) (row, x int) {
	rowStart := 0
	for {
		if n >= rowStart && n < rowStart+layout.rowLen(row) {
			x = n - rowStart
			if row%2 != 0 && layout.InvertOddRows {
				x = layout.rowLen(row) - 1 - x
			}
			return
		}
		rowStart += layout.rowLen(row)
		row++
	}
}

func hexagon(x, y, r float64) (out Poly) {
	rr := Sqrt3 / 2.0 * r
	out = make([]Point, 6)
	out[0] = Point{x, y - r}
	out[1] = Point{x + rr, y - r/2}
	out[2] = Point{x + rr, y + r/2}
	out[3] = Point{x, y + r}
	out[4] = Point{x - rr, y + r/2}
	out[5] = Point{x - rr, y - r/2}
	return
}

func (layout *HexLayout) Cell(n int) Poly {
	if n >= layout.NumCells() {
		return nil
	}

	row, xpos := layout.cellPos(n)

	rr := Size / float64(layout.EdgeLen*2-1) / 2.0
	r := rr * 2.0 / Sqrt3

	ystep := r + r/2.0
	xstep := rr * 2

	y0 := Size/2.0 - float64(layout.EdgeLen-1)*ystep
	y := y0 + float64(row)*ystep

	var x0 float64
	if row < layout.EdgeLen {
		x0 = float64(layout.EdgeLen-row) * rr
	} else {
		x0 = float64(row-layout.EdgeLen+2) * rr
	}
	x := x0 + xstep*float64(xpos)

	return hexagon(x, y, r)
}

func (layout *HexLayout) Width() float64 {
	return Size
}

func (layout *HexLayout) Height() float64 {
	return Size
}
