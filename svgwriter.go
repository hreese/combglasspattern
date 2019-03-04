package combglasspattern

import (
	"fmt"
	"github.com/ajstarks/svgo/float"
	"math"
	"os"
)

const (
	CenterMarkerlaenge           = 10
	CenterMarkerStyle            = "stroke:#000000;stroke-opacity:1;stroke-width:0.35277778;stroke-linecap:round"
	CenterMarkerAnnotationOffset = CenterMarkerlaenge / 1.8
	CenterMarkerAnnotationStyle  = "font-family:sans-serif;font-weight:normal;font-style:normal;font-stretch:normal;font-variant:normal;font-size:3.52777777px;text-anchor:middle;text-align:center;"
	InnerCircleStyle             = "fill:none;stroke:#000000;stroke-opacity:1;stroke-width:1;"
	OuterCircleStyle             = "fill:#bdbdbd;stroke:#8e8e8e;stroke-opacity:1;stroke-width:0.5;stroke-dasharray:1,2;fill-opacity:0.29019609;stroke-dashoffset:0"
	OriginMarkLength             = 8
	BoardLineStyle               = "stroke:#000000;stroke-opacity:1;stroke-width:0.3;fill:#DDDDDD;" + "fill-rule:evenodd;"
)

func SVGCenterMarker(canvas *svg.SVG, center Point, s ...string) {
	var x, y float64 = center.X, center.Y
	// mark center
	canvas.TranslateRotate(x, y, 45)
	canvas.Line(-CenterMarkerlaenge/2, 0, CenterMarkerlaenge/2, 0, s...)
	canvas.Line(0, CenterMarkerlaenge/2, 0, -CenterMarkerlaenge/2, s...)
	canvas.Gend()
	// add coordinates
	canvas.Text(x, y-CenterMarkerAnnotationOffset, fmt.Sprintf("→%.1fmm ↓%.1fmm", x, y), CenterMarkerAnnotationStyle)
}

func SVGThroughhole(canvas *svg.SVG, center Point, innerRadius, outerRadius float64, numsides int) {
	var x, y = center.X, center.Y

	canvas.Group()
	// round Glass
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
	SVGCenterMarker(canvas, center, CenterMarkerStyle)
	canvas.Gend()
}

func SVGDrawBoard(canvas *svg.SVG, board BoardConfiguration, markOrigin bool) {
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
		BoardLineStyle)
	if markOrigin {
		// mark markOrigin
		canvas.Polygon([]float64{0, OriginMarkLength, 0}, []float64{0, 0, OriginMarkLength}, "stroke:none;fill:#000000;fill-opacity:1")
	}
	canvas.Gend()
}


func WriteSVG(filenameprefix string, variant Variant) error {
	var (
		f        *os.File
		err      error
		canvas   *svg.SVG
		filename string
	)
	filename = fmt.Sprintf("%s.svg", filenameprefix)

	// open svg file
	f, err = os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	canvas = svg.New(f)
	canvas.Startunit(variant.Board.Width, variant.Board.Height, "mm",
		fmt.Sprintf(`viewBox="0 0 %f %f"`, variant.Board.Width, variant.Board.Height))
	SVGDrawBoard(canvas, variant.Board, true)
	canvas.Group("Holes")
	for _, p := range variant.Points {
		SVGThroughhole(canvas, p, variant.Glass.InnerRadius, variant.Glass.OuterRadius, variant.Glass.NumberOfSides)
	}
	canvas.Gend()
	canvas.End()

	return nil
}
