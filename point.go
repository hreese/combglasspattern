package combglasspattern

import "math"

// Point denotes a point in a cartesian 2D plane
type Point struct {
	X, Y float64
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

// BoundingBox returns the smallest box encompassing all points
func BoundingBox(points []Point) (Point, Point) {
	var (
		xmin, xmax, ymin, ymax float64
	)
	if len(points) < 1 {
		return Point{0, 0}, Point{0, 0}
	}

	xmin, xmax, ymin, ymax = points[0].X, points[0].X, points[0].Y, points[0].Y
	for _, p := range points {
		xmin = math.Min(xmin, p.X)
		xmax = math.Max(xmax, p.X)
		ymin = math.Min(ymin, p.Y)
		ymax = math.Max(ymax, p.Y)
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
	//HolesMidPoint = MidPoint(BoundingBox(points))
	BoardMidPoint = MidPoint(Point{0, 0}, Point{board.Width, board.Height})
	CorrectionVector = (BoardMidPoint.Minus(HolesMidPoint))

	//litter.Dump(HolesMidPoint, BoardMidPoint, CorrectionVector)

	centeredPoints := MovePoints(points, CorrectionVector)

	return centeredPoints
}

func GenerateHoles(board BoardConfiguration, glass GlassConfiguration) ([]Point, []Point, []Point) {
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
	// hole midpoints can be generated within these bounds
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

	//return HolesSquare, HolesHexOne, HolesHexTwo
	return CenterAllHoles(HolesSquare, board), CenterAllHoles(HolesHexOne, board), CenterAllHoles(HolesHexTwo, board)
}
