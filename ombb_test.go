package ombb

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/fogleman/gg"
	r3 "github.com/golang/geo/r3"
	"github.com/markus-wa/quickhull-go/v2"
	"github.com/wsw0108/concaveman-go"
	"github.com/wsw0108/sphericalmercator-go"
)

var pointsNormal = []Point{
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

func TestConvexHull(t *testing.T) {
	diff := Point{250, 1050}
	points := make([]Point, len(pointsNormal))
	copy(points, pointsNormal)
	for i := range points {
		points[i] = points[i].Mul(1.5)
		points[i] = points[i].Diff(diff)
	}
	hull := ConvexHull(points)
	expected := []Point{
		{165.5, -91.5},
		{270.5, -133.5},
		{350, 15},
		{285.5, 82.5},
		{102.5, 111},
		{95, 15},
	}
	if !reflect.DeepEqual(hull, expected) {
		t.Error("TestConvexHull")
	}
}

func TestOmbb(t *testing.T) {
	diff := Point{250, 1050}
	points := make([]Point, len(pointsNormal))
	copy(points, pointsNormal)
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
	delta := 1e-6
	for i := range obb {
		if !obb[i].AlmostEquals(expected[i], delta) {
			t.Fatalf("Ombb(), got %v, want %v", obb, expected)
		}
	}
}

var pointsLonLat = []Point{
	{114.26671390000001, 30.599383600000003},
	{114.2668615, 30.599415800000003},
	{114.2670039, 30.599465600000002},
	{114.26715410000001, 30.599522},
	{114.26728890000001, 30.5995807},
	{114.2674544, 30.599659300000003},
	{114.2676225, 30.599737800000003},
	{114.2677881, 30.599833},
	{114.26793830000001, 30.599927},
	{114.26796370000001, 30.599949300000002},
	{114.2679679, 30.5999763},
	{114.26796030000001, 30.6000012},
	{114.2679392, 30.6000149},
	{114.26791100000001, 30.6000225},
	{114.26788710000001, 30.600023600000004},
	{114.26785090000001, 30.600088000000003},
	{114.2678795, 30.6001009},
	{114.26785190000001, 30.6001501},
	{114.26782480000001, 30.6001395},
	{114.2677749, 30.600238100000002},
	{114.26705270000001, 30.600096200000003},
	{114.26706820000001, 30.600001400000004},
	{114.2671476, 30.599873700000003},
	{114.2671262, 30.5998631},
	{114.2671706, 30.599784500000002},
	{114.267088, 30.5997406},
	{114.2670116, 30.599696700000003},
	{114.26693920000001, 30.5997975},
	{114.2669137, 30.599772400000003},
	{114.2668821, 30.599753600000003},
	{114.2668529, 30.599738900000002},
	{114.26681520000001, 30.599723200000003},
	{114.2667763, 30.599713800000004},
	{114.2667375, 30.599710700000003},
	{114.26669980000001, 30.599710700000003},
	{114.2666597, 30.599714900000002},
	{114.26662320000001, 30.599724300000002},
	{114.26659040000001, 30.5997379},
	{114.26662470000001, 30.599676600000002},
	{114.26660480000001, 30.599668500000003},
	{114.26664950000001, 30.5995893},
	{114.26666990000001, 30.599595700000002},
	{114.2667155, 30.5995166},
	{114.26665700000001, 30.5994761},
	{114.26663330000001, 30.599446500000003},
	{114.26662820000001, 30.599412200000003},
	{114.26664550000001, 30.599387600000004},
	{114.2666743, 30.599377500000003},
	{114.26671390000001, 30.599383600000003},
}

var transformerMercator = func(p Point, inverse bool) Point {
	mercator := sphericalmercator.New()
	var v []float64
	if !inverse {
		v = mercator.Forward(p[:])
	} else {
		v = mercator.Inverse(p[:])
	}
	return Point{v[0], v[1]}
}

func TestOmbbForLonLat(t *testing.T) {
	opt := Options{
		Transform: transformerMercator,
	}
	obb := Ombb(pointsLonLat, opt)
	expected := [4]Point{
		{114.26639138410215, 30.599784771864584},
		{114.26669797832271, 30.599302420611224},
		{114.26801133406653, 30.599920910685253},
		{114.26770473984595, 30.600403258859338},
	}
	delta := 1e-6
	for i := range obb {
		if !obb[i].AlmostEquals(expected[i], delta) {
			t.Fatalf("OmbbForLonLat(), got %v, want %v", obb, expected)
		}
	}
}

func convexHullQuickHull(points []Point) []Point {
	var pointsR3 []r3.Vector
	for _, p := range points {
		pointsR3 = append(pointsR3, r3.Vector{X: p[0], Y: p[1], Z: 0})
	}
	hull := new(quickhull.QuickHull).ConvexHull(pointsR3, false, true, 0)
	var result []Point
	for _, v := range hull.Vertices {
		result = append(result, Point{v.X, v.Y})
	}
	return result
}

func convexHullConcaveman(points []Point) []Point {
	var points2 []concaveman.Point
	for _, p := range points {
		points2 = append(points2, concaveman.Point{p[0], p[1]})
	}
	hull := concaveman.Concaveman(points2)
	var result []Point
	for _, v := range hull {
		result = append(result, Point{v[0], v[1]})
	}
	return result
}

func TestConvexHullCompare(t *testing.T) {
	t.Skip()
	diff := Point{250, 1050}
	points := make([]Point, len(pointsNormal))
	copy(points, pointsNormal)
	for i := range points {
		points[i] = points[i].Mul(1.5)
		points[i] = points[i].Diff(diff)
	}
	expected := []Point{
		{165.5, -91.5},
		{270.5, -133.5},
		{350, 15},
		{285.5, 82.5},
		{102.5, 111},
		{95, 15},
	}
	type args struct {
		convexHullFn ConvexHullFunc
		points       []Point
	}
	tests := []struct {
		name string
		args args
		want []Point
	}{
		{
			name: "ConvexHull",
			args: args{convexHullFn: ConvexHull, points: points},
			want: expected,
		},
		{
			name: "QuickHull",
			args: args{convexHullFn: convexHullQuickHull, points: points},
			want: expected,
		},
		{
			name: "Concaveman",
			args: args{convexHullFn: convexHullConcaveman, points: points},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.convexHullFn(tt.args.points); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvexHull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDrawConvexHull(t *testing.T) {
	diff := Point{250, 1050}
	points := make([]Point, len(pointsNormal))
	copy(points, pointsNormal)
	for i := range points {
		points[i] = points[i].Mul(1.5)
		points[i] = points[i].Diff(diff)
	}
	off := Point{0, 200}
	drawLine := func(dc *gg.Context, p1, p2 Point) {
		dc.DrawLine(p1[0]+off[0], p1[1]+off[1], p2[0]+off[0], p2[1]+off[1])
	}
	output := func(dc *gg.Context, hull []Point, filename string) error {
		{
			dc.SetColor(color.RGBA{0, 0, 0, 255})
			p1 := points[len(points)-1]
			for _, p2 := range points {
				drawLine(dc, p1, p2)
				p1 = p2
			}
			dc.Stroke()
		}
		{
			dc.SetColor(color.RGBA{255, 0, 0, 255})
			p1 := hull[len(hull)-1]
			for _, p2 := range hull {
				drawLine(dc, p1, p2)
				p1 = p2
			}
			dc.Stroke()
		}
		return dc.SavePNG(filename)
	}
	{
		dc := gg.NewContext(800, 600)
		hull := ConvexHull(points)
		err := output(dc, hull, "output-hull-default.png")
		if err != nil {
			panic(err)
		}
	}
	{
		dc := gg.NewContext(800, 600)
		hull := convexHullQuickHull(points)
		err := output(dc, hull, "output-hull-quickhull.png")
		if err != nil {
			panic(err)
		}
	}
	{
		dc := gg.NewContext(800, 600)
		hull := convexHullConcaveman(points)
		err := output(dc, hull, "output-hull-concaveman.png")
		if err != nil {
			panic(err)
		}
	}
}
