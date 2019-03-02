package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/ajstarks/svgo/float"
	. "gitlab.com/hreese/combglasspattern"
)

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

func GenerateHolesNew(board BoardConfiguration, glass GlassConfiguration) ([]Point, []Point, []Point) {
	var (
		EdgeOffset, GlassOffset, OffsetOne, OffsetTwo float64
		UpperLeft, LowerRight                         Point
		HolesSquare, HolesHexOne, HolesHexTwo         []Point
	)
	// minimal distance from board's edge
	EdgeOffset = math.Max(
		board.WallOffset+glass.InnerRadius,
		glass.OuterRadius)
	// minimal distance between two glasses
	GlassOffset = math.Max(
		2*glass.InnerRadius+board.MinHoleDistance,
		2*glass.OuterRadius)
	// bounding box for glass midpoints
	UpperLeft = Point{EdgeOffset, EdgeOffset}
	LowerRight = Point{board.Width - EdgeOffset, board.Height - EdgeOffset}

	// variant "square"
	for y := UpperLeft.Y; y <= LowerRight.Y; y += GlassOffset {
		for x := UpperLeft.X; x <= LowerRight.X; x += GlassOffset {
			HolesSquare = append(HolesSquare, Point{x, y})
		}
	}

	// variants hex
	odd := true
	for y := UpperLeft.Y; y <= LowerRight.Y; y += GlassOffset * math.Sin(math.Pi/3.0) {
		// alternate x offsets
		if odd {
			OffsetOne, OffsetTwo = 0, GlassOffset/2
		} else {
			OffsetOne, OffsetTwo = GlassOffset/2, 0
		}
		// one
		for x := UpperLeft.X + OffsetOne; x <= LowerRight.X; x += GlassOffset {
			HolesHexOne = append(HolesHexOne, Point{x, y})
		}
		// two
		for x := UpperLeft.X + OffsetTwo; x <= LowerRight.X; x += GlassOffset {
			HolesHexTwo = append(HolesHexTwo, Point{x, y})
		}
		odd = !odd
	}

	return HolesSquare, HolesHexOne, HolesHexTwo
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

type Variant struct {
	points []Point
	board  BoardConfiguration
	glass  GlassConfiguration
}

func main() {
	var (
		f        *os.File
		err      error
		canvas   *svg.SVG
		variants = make(map[string]Variant)
		//board    = PresetsBoard["DadantWeber"]
		//glass    = PresetGlas["BienenRuckWabengläserRund500"]
		board = PresetsBoard["ZanderSpec"]
		glass = PresetGlas["TestGlas"]
	)

	square, hexOne, hexTwo := GenerateHolesNew(board, glass)
	variants[`Square.svg`] = Variant{square, board, glass}
	variants[`HexOne.svg`] = Variant{hexOne, board, glass}
	variants[`HexTwo.svg`] = Variant{hexTwo, board, glass}
	fmt.Printf("Sqare:   %d\nHexOne:  %d\nHexTwo:  %d\n", len(square), len(hexOne), len(hexTwo))
	if !board.IsSquare() {
		turnedBoard := board.Turn90()
		_, hexThree, hexFour := GenerateHolesNew(turnedBoard, glass)
		variants[`HexThree.svg`] =  Variant{hexThree, turnedBoard, glass}
		variants[`HexFour.svg`] =  Variant{hexFour, turnedBoard, glass}
		fmt.Printf("HexThree:  %d\nHexFour:  %d\n", len(hexThree), len(hexFour))
	}

	for filename, variant := range variants {
		f, err = os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		canvas = svg.New(f)
		canvas.Startunit(variant.board.Width, variant.board.Height, "mm",
			fmt.Sprintf(`viewBox="0 0 %f %f"`, variant.board.Width, variant.board.Height))
		DrawBoard(canvas, variant.board, true)
		canvas.Group("Holes")
		for _, p := range variant.points {
			Throughhole(canvas, p, variant.glass.InnerRadius, variant.glass.OuterRadius, variant.glass.NumberOfSides)
		}
		canvas.Gend()
		canvas.End()
	}

}
