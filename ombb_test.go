package ombb

import (
	"testing"
)

func TestOmbb(t *testing.T) {
	points := []Point{
		{235, 774},
		{245, 740},
		{230, 710},
		{240, 703},
		{274, 733},
		{306, 710},
		{272, 690},
		{277, 639},
		{305, 645},
		{347, 611},
		{340, 639},
		{298, 674},
		{325, 702},
		{335, 663},
		{355, 645},
		{350, 686},
		{400, 710},
		{360, 725},
		{357, 755},
		{328, 723},
		{291, 741},
		{289, 754},
		{266, 757},
	}
	diff := Point{250, 1050}
	for i := range points {
		points[i] = points[i].Mul(1.5)
		points[i] = points[i].Diff(diff)
	}
	obb := Ombb(points)
	expected := [4]Point{
		{41.36206896551724, -41.84482758620689},
		{287.82758620689657, -140.4310344827586},
		{364.37931034482756, 50.94827586206896},
		{117.91379310344828, 149.5344827586207},
	}
	for i := range obb {
		if !obb[i].AlmostEquals(expected[i], 1e-9) {
			t.Fatalf("Ombb()[%d], got %v, want %v", i, obb[i], expected[i])
		}
	}
}
