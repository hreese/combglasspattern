package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/sanity-io/litter"

	"github.com/ajstarks/svgo/float"
)

type Point struct {
	X, Y float64
}

type BoardConfiguration struct {
	Width, Height   float64
	WallOffset      float64
	MinHoleDistance float64
}

type GlassConfiguration struct {
	InnerRadius   float64
	OuterRadius   float64
	NumberOfSides int
}

var (
	DadantWeber BoardConfiguration = BoardConfiguration{
		Width:           490,
		Height:          490,
		WallOffset:      10,
		MinHoleDistance: 10,
	}
	TestBrett BoardConfiguration = BoardConfiguration{
		Width:           580,
		Height:          510,
		WallOffset:      10,
		MinHoleDistance: 10,
	}

	// https://www.holtermann-glasshop.de/Sechseckglaeser/Sechseckglas-580-ml/
	HolterMannTwistOffSechseckglas580 GlassConfiguration = GlassConfiguration{
		InnerRadius:   82 / 2,
		OuterRadius:   95 / 2,
		NumberOfSides: 6,
	}
)

const (
	CenterMarkerlaenge           = 10
	CenterMarkerStyle            = "stroke:#000000;stroke-opacity:1;stroke-width:0.35277778;stroke-linecap:round"
	CenterMarkerAnnotationOffset = CenterMarkerlaenge / 1.8
	CenterMarkerAnnotationStyle  = "font-family:sans-serif;font-weight:normal;font-style:normal;font-stretch:normal;font-variant:normal;font-size:3.52777777px;text-anchor:middle;text-align:center;"
	InnerCircleStyle             = "fill:none;stroke:#000000;stroke-opacity:1;stroke-width:1;"
	OuterCircleStyle             = "fill:#bdbdbd;stroke:#8e8e8e;stroke-opacity:1;stroke-width:1;stroke-dasharray:1,2;fill-opacity:0.29019609;stroke-dashoffset:0"
)

func CenterMarker(canvas *svg.SVG, x, y float64) {
	// mark center
	canvas.TranslateRotate(x, y, 45)
	canvas.Line(-CenterMarkerlaenge/2, 0, CenterMarkerlaenge/2, 0, CenterMarkerStyle)
	canvas.Line(0, CenterMarkerlaenge/2, 0, -CenterMarkerlaenge/2, CenterMarkerStyle)
	canvas.Gend()
	// add coordinates
	canvas.Text(x, y-CenterMarkerAnnotationOffset, fmt.Sprintf("→%.1fmm ↓%.1fmm", x, y), CenterMarkerAnnotationStyle)
}

func Throughhole(canvas *svg.SVG, x, y, innerRadius, outerRadius float64, numsides int) {
	// assume a round glass
	if numsides < 3 {
		canvas.Circle(x, y, outerRadius, OuterCircleStyle)
	} else {
		canvas.TranslateRotate(x, y, 0)
		xcoords := make([]float64, numsides)
		ycoords := make([]float64, numsides)
		angleStep := math.Pi * 2 / float64(numsides)
		for step := 0; step < numsides; step++ {
			xcoords[step] = outerRadius * math.Cos(float64(step)*angleStep)
			ycoords[step] = outerRadius * math.Sin(float64(step)*angleStep)
		}
		canvas.Polygon(xcoords, ycoords, OuterCircleStyle)
		canvas.Gend()
	}
	canvas.Circle(x, y, innerRadius, InnerCircleStyle)
	canvas.Text(x, y+innerRadius/2, fmt.Sprintf("Ø %.1fmm", innerRadius*2), CenterMarkerAnnotationStyle)
	CenterMarker(canvas, x, y)

}

func CenterAllHoles(points []Point, board BoardConfiguration) []Point {
	var (
		xmin, xmax, ymin, ymax float64
	)
	if len(points) < 1 {
		return points
	}

	// find bounding box
	xmin, xmax, ymin, ymax = points[0].X, points[0].X, points[0].Y, points[0].Y
	for _, p := range points {
		if p.X < xmin {
			xmin = p.X
		}
		if p.X > xmax {
			xmax = p.X
		}
		if p.Y < ymin {
			ymin = p.Y
		}
		if p.Y > ymax {
			ymax = p.Y
		}
	}

	// TODO: fix
	// calculate offset by comparing the middle points
	xoff := (board.Width / 2) - (xmax-xmin)/2
	yoff := (board.Height / 2) - (ymax-ymin)/2
	litter.Dump(board, xmin, xmax, ymin, ymax, xoff, yoff)

	for idx := range points {
		points[idx].X = points[idx].X - xoff
		points[idx].Y = points[idx].Y - yoff
	}
	return points
}

func GenerateBoard(board BoardConfiguration, glass GlassConfiguration) ([]Point, []Point) {
	var (
		holeDistance = 2*glass.InnerRadius + board.MinHoleDistance
		sideOffset   = board.WallOffset + glass.InnerRadius
		squareholes  = make([]Point, 0)
		hexholes     = make([]Point, 0)
		xoff         float64
	)

	// find minimal hole distance
	if 2*glass.OuterRadius > holeDistance {
		holeDistance = 2 * glass.OuterRadius
	}

	// find minimal initial wall distance
	if glass.OuterRadius > sideOffset {
		sideOffset = glass.OuterRadius
	}

	// variant "square"
	for x := sideOffset; x < board.Width-board.WallOffset-glass.InnerRadius; x += holeDistance {
		for y := sideOffset; y < board.Height-board.WallOffset-glass.InnerRadius; y += holeDistance {
			squareholes = append(squareholes, Point{x, y})
		}
	}

	// variant "hex"
	odd := true
	for y := sideOffset; y < board.Height-board.WallOffset-glass.InnerRadius; y += holeDistance * math.Sin(math.Pi/3.0) {
		if !odd {
			xoff = holeDistance / 2
		} else {
			xoff = 0
		}
		for x := xoff + sideOffset; x < board.Width-board.WallOffset-glass.InnerRadius; x += holeDistance {
			hexholes = append(hexholes, Point{x, y})
		}
		odd = !odd
	}

	//return CenterAllHoles(squareholes, board), CenterAllHoles(hexholes, board)
	return squareholes, hexholes
}

func main() {
	var (
		f      *os.File
		err    error
		canvas *svg.SVG
		board  = TestBrett
		glass  = HolterMannTwistOffSechseckglas580
	)

	square, hex := GenerateBoard(board, glass)
	fmt.Printf("Sqare: %d\nHex:   %d\n", len(square), len(hex))

	f, err = os.OpenFile("square.svg", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	canvas = svg.New(f)
	canvas.Startunit(board.Width, board.Height, "mm", fmt.Sprintf(`viewBox="0 0 %f %f"`, board.Width, board.Height))
	for _, p := range square {
		Throughhole(canvas, p.X, p.Y, glass.InnerRadius, glass.OuterRadius, glass.NumberOfSides)
	}
	canvas.End()
	f.Close()

	f, err = os.OpenFile("hex.svg", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	canvas = svg.New(f)
	canvas.Startunit(board.Width, board.Height, "mm", fmt.Sprintf(`viewBox="0 0 %f %f"`, board.Width, board.Height))
	for _, p := range hex {
		Throughhole(canvas, p.X, p.Y, glass.InnerRadius, glass.OuterRadius, glass.NumberOfSides)
	}
	canvas.End()
	f.Close()

}
