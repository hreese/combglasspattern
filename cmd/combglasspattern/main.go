package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/ajstarks/svgo/float"
	. "gitlab.com/hreese/combglasspattern"
)

const (
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

func main() {
	var (
		f        *os.File
		err      error
		canvas   *svg.SVG
		variants = make(map[string][]Point)
		//board    = PresetsBoard["DadantWeber"]
		//glass    = PresetGlas["BienenRuckWabengläserRund500"]
		board    = PresetsBoard["DemoBrettA4"]
		glass    = PresetGlas["DemoGlasEckig"]
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
