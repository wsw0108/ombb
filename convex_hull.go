package ombb

import "math"

const (
	on    = 0
	left  = 1
	right = 2
)

const (
	defaultAlmostZero = 0.00001
)

func getSideOfLine(lineStart, lineEnd, point Point, almostZero float64) int {
	d := (lineEnd[0]-lineStart[0])*(point[1]-lineStart[1]) - (lineEnd[1]-lineStart[1])*(point[0]-lineStart[0])
	switch {
	case d > almostZero:
		return left
	case d < -almostZero:
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

func ConvexHull(points []Point, precision ...float64) []Point {
	if len(points) < 3 {
		return points
	}

	almostZero := defaultAlmostZero
	if len(precision) > 0 {
		almostZero = precision[0]
	}

	hullPt := points[0]
	var convexHull []Point

	for _, p := range points {
		if p[0] < hullPt[0] {
			hullPt = p
		} else if math.Abs(p[0]-hullPt[0]) < almostZero {
			if p[1] < hullPt[1] {
				hullPt = p
			}
		}
	}

	endPt := points[0]
	for {
		convexHull = append(convexHull, hullPt)

		for _, p := range points[1:] {
			side := getSideOfLine(hullPt, endPt, p, almostZero)

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
