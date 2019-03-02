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