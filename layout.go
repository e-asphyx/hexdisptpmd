package main

import "math"

type Point struct {
	X float64
	Y float64
}

type Poly []Point

type Layout interface {
	Width() float64
	Height() float64
	NumCells() int
	Cell(n int) Poly
}

func (points Poly) BoundingBox() (tl, br Point) {
	if len(points) == 0 {
		return
	}

	xmin := points[0].X
	xmax := points[0].X

	ymin := points[0].Y
	ymax := points[0].Y

	for _, p := range points {
		xmin = math.Min(xmin, p.X)
		xmax = math.Max(xmax, p.X)

		ymin = math.Min(ymin, p.Y)
		ymax = math.Max(ymax, p.Y)
	}
	return Point{xmin, ymin}, Point{xmax, ymax}
}

func (poly Poly) PointInside(point Point) bool {
	if len(poly) < 2 {
		return false
	}

	cnt := 0
	for i := range poly {
		begin := poly[i]
		end := poly[(i+1)%len(poly)]

		vx := end.X - begin.X
		vy := end.Y - begin.Y
		vpx := point.X - begin.X
		vpy := point.Y - begin.Y

		prod := vpx*vy - vpy*vx
		if prod <= 0 {
			cnt++
		}
	}
	return cnt == len(poly)
}
