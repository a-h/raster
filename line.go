package raster

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// Line defines a line between two points in 2D space.
type Line struct {
	From   image.Point
	To     image.Point
	points map[image.Point]interface{}
}

// NewLine creates a new line between the specified points and precalculates
// the points which the line passes through.
func NewLine(fromX, fromY int, toX, toY int) *Line {
	l := &Line{
		From:   image.Point{fromX, fromY},
		To:     image.Point{toX, toY},
		points: make(map[image.Point]interface{}),
	}
	l.calculatePoints()
	return l
}

// String provides a string representation of the line in form "{x, y} to {x, y}"
func (l *Line) String() string {
	return fmt.Sprintf("{%v, %v} to {%v, %v}", l.From.X, l.From.Y, l.To.X, l.To.Y)
}

func (l *Line) calculatePoints() {
	for _, p := range line(l.From.X, l.From.Y, l.To.X, l.To.Y) {
		l.points[p] = true
	}
}

// Points returns the precalculated list of points which the line will pass through.
func (l *Line) Points() (points []image.Point) {
	for k := range l.points {
		points = append(points, k)
	}
	return points
}

// ContainsPoint returns true if a point appears on the line.
func (l *Line) ContainsPoint(p image.Point) bool {
	_, ok := l.points[p]
	return ok
}

// Draw draws out the line onto the provided image, using the specified colour. It also
// returns the points which were drawn.
func (l *Line) Draw(img *image.RGBA, c color.RGBA) (points []image.Point) {
	points = l.Points()
	for _, p := range points {
		img.Set(p.X, p.Y, c)
	}
	return points
}

func line(fromX, fromY int, toX, toY int) (points []image.Point) {
	// Vertical line.
	if fromX == toX {
		if toY < fromY {
			toX, toY, fromX, fromY = fromX, fromY, toX, toY
		}
		for y := fromY; y <= toY; y++ {
			points = append(points, image.Point{fromX, y})
		}
		return
	}

	// We're moving from fromX to toX, so make sure they're in the right order.
	if toX < fromX {
		toX, toY, fromX, fromY = fromX, fromY, toX, toY
	}

	// Horizontal line, we don't need floating points.
	if fromY == toY {
		for x := fromX; x <= toX; x++ {
			points = append(points, image.Point{x, fromY})
		}
		return
	}

	// It's a slope.
	rise := toY - fromY
	run := toX - fromX

	xdelta := float64(run) / float64(rise)
	ydelta := float64(rise) / float64(run)

	if math.Abs(xdelta) < math.Abs(ydelta) {
		// We're moving from fromY to toY, so make sure they're in the right order.
		if toY < fromY {
			toX, toY, fromX, fromY = fromX, fromY, toX, toY
		}

		x := float64(fromX)
		for y := fromY; y < toY; y++ {
			points = append(points, image.Point{int(x), y})
			x += xdelta
		}
	} else {
		y := float64(fromY)
		for x := fromX; x < toX; x++ {
			points = append(points, image.Point{x, int(y)})
			y += ydelta
		}
	}
	points = append(points, image.Point{toX, toY})
	return
}
