package main

import (
	"fmt"
	"log"
	"math"
	"os"

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
	// https://www.holtermann-glasshop.de/Sechseckglaeser/Sechseckglas-580-ml/
	HolterMannTwistOffSechseckglas580 GlassConfiguration = GlassConfiguration{
		InnerRadius:   82 / 2,
		OuterRadius:   95 / 2,
		NumberOfSides: 6,
	}
	// https://www.holtermann-glasshop.de/Designglaeser/Viereckglas-312-ml/Viereckglas-312-ml-Biene.html
	HolterMannTwistOffViereckglas312 GlassConfiguration = GlassConfiguration{
		InnerRadius:   60 / 2,
		OuterRadius:   75 / 2,
		NumberOfSides: 4,
	}
	// https://www.flaschenbauer.de/einmachglaeser/sechskantglaeser/sechskantglas-580-ml-to-82
	FlaschenBauerSechskantglas580mlTO82 GlassConfiguration = GlassConfiguration{
		InnerRadius:   82 / 2,
		OuterRadius:   95 / 2,
		NumberOfSides: 6,
	}
	// https://www.bienen-ruck.de/imkershop/honigverkauf-werbemittel/twist-off-glaeser/1902/wabenglaeser-rund
	BienenRuckWabengläserRund500 GlassConfiguration = GlassConfiguration{
		InnerRadius:   82 / 2,
		OuterRadius:   90 / 2,
		NumberOfSides: 0,
	}
	TestBrett BoardConfiguration = BoardConfiguration{
		Width:           500,
		Height:          600,
		WallOffset:      10,
		MinHoleDistance: 10,
	}
	TestGlas GlassConfiguration = GlassConfiguration{
		InnerRadius:   60 / 2,
		OuterRadius:   90 / 2,
		NumberOfSides: 0,
	}
)

const (
	OriginMarkLength             = 8
	CenterMarkerlaenge           = 10
	CenterMarkerStyle            = "stroke:#000000;stroke-opacity:1;stroke-width:0.35277778;stroke-linecap:round"
	CenterMarkerAnnotationOffset = CenterMarkerlaenge / 1.8
	CenterMarkerAnnotationStyle  = "font-family:sans-serif;font-weight:normal;font-style:normal;font-stretch:normal;font-variant:normal;font-size:3.52777777px;text-anchor:middle;text-align:center;"
	InnerCircleStyle             = "fill:none;stroke:#000000;stroke-opacity:1;stroke-width:1;"
	OuterCircleStyle             = "fill:#bdbdbd;stroke:#8e8e8e;stroke-opacity:1;stroke-width:0.5;stroke-dasharray:1,2;fill-opacity:0.29019609;stroke-dashoffset:0"
)

func (board BoardConfiguration) CenterPoint() Point {
	return Point{board.Width / 2, board.Height / 2}
}

func (a Point) Plus(b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y}
}

func (a Point) Minus(b Point) Point {
	return Point{a.X - b.X, a.Y - b.Y}
}

func (a Point) Scale(fac float64) Point {
	return Point{a.X * fac, a.Y * fac}
}

func CenterMarker(canvas *svg.SVG, center Point, s ...string) {
	var x, y float64 = center.X, center.Y
	// mark center
	canvas.TranslateRotate(x, y, 45)
	canvas.Line(-CenterMarkerlaenge/2, 0, CenterMarkerlaenge/2, 0, s...)
	canvas.Line(0, CenterMarkerlaenge/2, 0, -CenterMarkerlaenge/2, s...)
	canvas.Gend()
	// add coordinates
	canvas.Text(x, y-CenterMarkerAnnotationOffset, fmt.Sprintf("→%.1fmm ↓%.1fmm", x, y), CenterMarkerAnnotationStyle)
}

func Throughhole(canvas *svg.SVG, center Point, innerRadius, outerRadius float64, numsides int) {
	var x, y = center.X, center.Y

	canvas.Group()
	// round glass
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
	CenterMarker(canvas, center, CenterMarkerStyle)
	canvas.Gend()
}

func BoundingBox(points []Point) (Point, Point) {
	var (
		xmin, xmax, ymin, ymax float64
	)
	if len(points) < 1 {
		return Point{0, 0}, Point{0, 0}
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
	return Point{xmin, ymin}, Point{xmax, ymax}
}

func Centroid(points []Point) Point {
	var c = Point{0, 0}
	for _, p := range points {
		c = c.Plus(p)
	}
	return c.Scale(1.0 / float64(len(points)))
}

func MovePoints(points []Point, vec Point) []Point {
	var newPoints []Point
	for idx := range points {
		newPoints = append(newPoints, points[idx].Plus(vec))
	}
	return newPoints
}

func MidPoint(a, b Point) Point {
	var (
		xmin, xmax, ymin, ymax float64 = a.X, b.X, a.Y, b.Y
	)
	if xmin > xmax {
		xmin, xmax = xmax, xmin
	}
	if ymin > ymax {
		ymin, ymax = ymax, ymin
	}

	return Point{(xmax - xmin) / 2, (ymax - ymin) / 2}
}

func CenterAllHoles(points []Point, board BoardConfiguration) []Point {
	var (
		HolesMidPoint    Point
		BoardMidPoint    Point
		CorrectionVector Point
	)
	if len(points) < 1 {
		return points
	}

	HolesMidPoint = Centroid(points)
	BoardMidPoint = MidPoint(Point{0, 0}, Point{board.Width, board.Height})
	CorrectionVector = (BoardMidPoint.Minus(HolesMidPoint))

	centeredPoints := MovePoints(points, CorrectionVector)

	return centeredPoints
}

func GenerateHoles(board BoardConfiguration, glass GlassConfiguration) ([]Point, []Point) {
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

	return CenterAllHoles(squareholes, board), CenterAllHoles(hexholes, board)
	//return squareholes, hexholes
}

func DrawBoard(canvas *svg.SVG, board BoardConfiguration, origin bool) {
	canvas.Group("Board")
	canvas.Path(fmt.Sprintf("M %f %f L %f %f L %f %f L %f %f L %f %f z M %f %f L %f %f L %f %f L %f %f L %f %f z",
		0.0, 0.0,
		board.Width, 0.0,
		board.Width, board.Height,
		0.0, board.Height,
		0.0, 0.0,
		board.WallOffset, board.WallOffset,
		board.Width-board.WallOffset, board.WallOffset,
		board.Width-board.WallOffset, board.Height-board.WallOffset,
		board.WallOffset, board.Height-board.WallOffset,
		board.WallOffset, board.WallOffset),
		"stroke:#000000;stroke-opacity:1;stroke-width:0.3;fill:#DDDDDD;"+"fill-rule:evenodd;")
	if origin {
		// mark origin
		canvas.Polygon([]float64{0, OriginMarkLength, 0}, []float64{0, 0, OriginMarkLength}, "stroke:none;fill:#000000;fill-opacity:1")
	}
	canvas.Gend()
}

func main() {
	var (
		f        *os.File
		err      error
		canvas   *svg.SVG
		variants = make(map[string][]Point)
		board    = DadantWeber
		glass    = BienenRuckWabengläserRund500
	)

	square, hex := GenerateHoles(board, glass)
	variants[`square.svg`] = square
	variants[`hex.svg`] = hex
	fmt.Printf("Sqare: %d\nHex:   %d\n", len(square), len(hex))

	for filename, points := range variants {
		f, err = os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		canvas = svg.New(f)
		canvas.Startunit(board.Width, board.Height, "mm", fmt.Sprintf(`viewBox="0 0 %f %f"`, board.Width, board.Height))
		DrawBoard(canvas, board, true)
		canvas.Group("Holes")
		for _, p := range points {
			Throughhole(canvas, p, glass.InnerRadius, glass.OuterRadius, glass.NumberOfSides)
		}
		canvas.Gend()
		canvas.End()
	}

}
