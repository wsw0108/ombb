package ombb

import "math"

const (
	on          = 0
	left        = 1
	right       = 2
	almost_zero = 0.00001
)

func getSideOfLine(lineStart, lineEnd, point Point) int {
	d := (lineEnd[0]-lineStart[0])*(point[1]-lineStart[1]) - (lineEnd[1]-lineStart[1])*(point[0]-lineStart[0])
	switch {
	case d > almost_zero:
		return left
	case d < -almost_zero:
		return right
	default:
		return on
	}
}

func reversePoints(points []Point) []Point {
	n := len(points)
	result := make([]Point, n)
	for i, p := range points {
		result[n-i-1] = p
	}
	return result
}

func ConvexHull(points []Point) []Point {
	if len(points) < 3 {
		return points
	}

	hullPt := points[0]
	var convexHull []Point

	for _, p := range points {
		if p[0] < hullPt[0] {
			hullPt = p
		} else if math.Abs(p[0]-hullPt[0]) < almost_zero {
			if p[1] < hullPt[1] {
				hullPt = p
			}
		}
	}

	endPt := points[0]
	for {
		convexHull = append(convexHull, hullPt)

		for _, p := range points[1:] {
			side := getSideOfLine(hullPt, endPt, p)

			// in case point lies on line take the one further away.
			// this fixes the collinearity problem.
			if endPt.Equals(hullPt) || (side == left || (side == on && hullPt.Distance(p) > hullPt.Distance(endPt))) {
				endPt = p
			}
		}

		hullPt = endPt

		if endPt.Equals(convexHull[0]) {
			break
		}
	}

	return reversePoints(convexHull)
}
