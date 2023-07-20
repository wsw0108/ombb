package ombb

import "math"

type ombbContext struct {
	bestObb     [4]Point
	bestObbArea float64
}

type Options struct {
	ConvexHull func([]Point) []Point
}

func Ombb(points []Point, opts ...Options) [4]Point {
	opt := Options{
		ConvexHull: ConvexHull,
	}
	if len(opts) > 0 {
		opt = opts[0]
	}
	convexHull := opt.ConvexHull(points)
	ctx := &ombbContext{}
	return ctx.calcOmbb(convexHull)
}

func intersectLines(start0, dir0, start1, dir1 Point) Point {
	dd := dir0.Cross(dir1)
	// dd=0 => lines are parallel. we don't care as our lines are never parallel.
	d := start1.Diff(start0)
	t := d.Cross(dir1) / dd
	return Point{
		start0[0] + t*dir0[0],
		start0[1] + t*dir0[1],
	}
}

func (ctx *ombbContext) updateOmbb(leftStart, leftDir, rightStart, rightDir, topStart, topDir, bottomStart, bottomDir Point) {
	obbUpperLeft := intersectLines(leftStart, leftDir, topStart, topDir)
	obbUpperRight := intersectLines(rightStart, rightDir, topStart, topDir)
	obbBottomLeft := intersectLines(bottomStart, bottomDir, leftStart, leftDir)
	obbBottomRight := intersectLines(bottomStart, bottomDir, rightStart, rightDir)
	distLeftRight := obbUpperLeft.Distance(obbUpperRight)
	distTopBottom := obbUpperLeft.Distance(obbBottomLeft)
	obbArea := distLeftRight * distTopBottom

	if obbArea < ctx.bestObbArea {
		ctx.bestObb[0] = obbUpperLeft
		ctx.bestObb[1] = obbBottomLeft
		ctx.bestObb[2] = obbBottomRight
		ctx.bestObb[3] = obbUpperRight
		ctx.bestObbArea = obbArea
	}
}

func (ctx *ombbContext) calcOmbb(convexHull []Point) [4]Point {
	ctx.bestObbArea = math.MaxFloat64

	edgeDirs := make([]Point, len(convexHull))
	for i := range convexHull {
		p1 := convexHull[i]
		p2 := convexHull[(i+1)%len(convexHull)]
		edgeDirs[i] = p2.Diff(p1).Normalize()
	}

	minPt := Point{math.MaxFloat64, math.MaxFloat64}
	maxPt := Point{-math.MaxFloat64, -math.MaxFloat64}
	var leftIdx, rightIdx, topIdx, bottomIdx int

	for i, pt := range convexHull {
		if pt[0] < minPt[0] {
			minPt[0] = pt[0]
			leftIdx = i
		}

		if pt[0] > maxPt[0] {
			maxPt[0] = pt[0]
			rightIdx = i
		}

		if pt[1] < minPt[1] {
			minPt[1] = pt[1]
			bottomIdx = i
		}

		if pt[1] > maxPt[1] {
			maxPt[1] = pt[1]
			topIdx = i
		}
	}

	leftDir := Point{0, -1}
	rightDir := Point{0, 1}
	topDir := Point{-1, 0}
	bottomDir := Point{1, 0}

	for range convexHull {
		phis := []float64{
			math.Acos(leftDir.Dot(edgeDirs[leftIdx])),
			math.Acos(rightDir.Dot(edgeDirs[rightIdx])),
			math.Acos(topDir.Dot(edgeDirs[topIdx])),
			math.Acos(bottomDir.Dot(edgeDirs[bottomIdx])),
		}
		var minPhi float64 = math.MaxFloat64
		var lineWithSmallestAngle int
		for i, phi := range phis {
			if phi < minPhi {
				minPhi = phi
				lineWithSmallestAngle = i
			}
		}
		switch lineWithSmallestAngle {
		case 0:
			leftDir = edgeDirs[leftIdx]
			rightDir = leftDir.Negate()
			topDir = leftDir.Orthogonal()
			bottomDir = topDir.Negate()
			leftIdx = (leftIdx + 1) % len(convexHull)
		case 1:
			rightDir = edgeDirs[rightIdx]
			leftDir = rightDir.Negate()
			topDir = leftDir.Orthogonal()
			bottomDir = topDir.Negate()
			rightIdx = (rightIdx + 1) % len(convexHull)
		case 2:
			topDir = edgeDirs[topIdx]
			bottomDir = topDir.Negate()
			leftDir = bottomDir.Orthogonal()
			rightDir = leftDir.Negate()
			topIdx = (topIdx + 1) % len(convexHull)
		case 3:
			bottomDir = edgeDirs[bottomIdx]
			topDir = bottomDir.Negate()
			leftDir = bottomDir.Orthogonal()
			rightDir = leftDir.Negate()
			bottomIdx = (bottomIdx + 1) % len(convexHull)
		}

		ctx.updateOmbb(convexHull[leftIdx], leftDir, convexHull[rightIdx], rightDir, convexHull[topIdx], topDir, convexHull[bottomIdx], bottomDir)
	}

	return ctx.bestObb
}
