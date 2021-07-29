package ombb

import "math"

type Point [2]float64

func (p Point) Mul(m float64) Point {
	return Point{m * p[0], m * p[1]}
}

func (p Point) Normalize() Point {
	d := p.Length()
	return Point{p[0] / d, p[1] / d}
}

func (p *Point) NormalizeInPlace() {
	d := p.Length()
	p[0] /= d
	p[1] /= d
}

func (p Point) Negate() Point {
	return Point{-p[0], -p[1]}
}

func (p Point) Length() float64 {
	return math.Hypot(p[0], p[1])
}

func (p Point) Diff(op Point) Point {
	return Point{p[0] - op[0], p[1] - op[1]}
}

func (p Point) Distance(op Point) float64 {
	return math.Hypot(p[0]-op[0], p[1]-op[1])
}

func (p Point) Dot(op Point) float64 {
	return p[0]*op[0] + p[1]*op[1]
}

func (p Point) Cross(op Point) float64 {
	return p[0]*op[1] - p[1]*op[0]
}

func (p Point) Equals(op Point) bool {
	return p[0] == op[0] && p[1] == op[1]
}

func almostEqual(f1, f2, delta float64) bool {
	return math.Abs(f1-f2) <= delta
}

func (p Point) AlmostEquals(op Point, delta float64) bool {
	return almostEqual(p[0], op[0], delta) && almostEqual(p[1], op[1], delta)
}

func (p Point) Orthogonal() Point {
	return Point{p[1], -p[0]}
}
