package raster

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Line defines a line between two points in 2D space.
type Line struct {
	From         image.Point
	To           image.Point
	OutlineColor color.RGBA
}

// NewLine creates a new line between the specified points.
func NewLine(from image.Point, to image.Point, outlineColor color.RGBA) *Line {
	l := &Line{
		From:         from,
		To:           to,
		OutlineColor: outlineColor,
	}
	return l
}

// String provides a string representation of the line in form "{x, y} to {x, y}"
func (l *Line) String() string {
	return fmt.Sprintf("{%v, %v} to {%v, %v}", l.From.X, l.From.Y, l.To.X, l.To.Y)
}

// Points returns the points of the line will pass through.
func (l *Line) Points() (points []image.Point) {
	accumulator := func(x, y int) bool {
		points = append(points, image.Point{x, y})
		return true
	}
	line(l.From.X, l.From.Y, l.To.X, l.To.Y, accumulator)
	return points
}

// ContainsPoint returns true if a point appears on the line.
func (l *Line) ContainsPoint(p image.Point) bool {
	contains := false
	containerCheck := func(x, y int) bool {
		if p.X == x && p.Y == y {
			contains = true
		}
		return !contains
	}
	line(l.From.X, l.From.Y, l.To.X, l.To.Y, containerCheck)
	return contains
}

// Draw draws the element to the img, img could be an image.RGBA* or screen buffer.
func (l *Line) Draw(img draw.Image) {
	drawer := func(x, y int) bool {
		img.Set(x, y, l.OutlineColor)
		return true
	}
	line(l.From.X, l.From.Y, l.To.X, l.To.Y, drawer)
}

func line(fromX, fromY int, toX, toY int, f func(x, y int) bool) {
	// Vertical line.
	if fromX == toX {
		if toY < fromY {
			toX, toY, fromX, fromY = fromX, fromY, toX, toY
		}
		for y := fromY; y <= toY; y++ {
			if !f(fromX, y) {
				return
			}
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
			if !f(x, fromY) {
				return
			}
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
			if !f(int(x), y) {
				return
			}
			x += xdelta
		}
	} else {
		y := float64(fromY)
		for x := fromX; x < toX; x++ {
			if !f(x, int(y)) {
				return
			}
			y += ydelta
		}
	}
	f(toX, toY)
}

// Bounds provides the area of the bounding box of the line.
func (l *Line) Bounds() image.Rectangle {
	first := true
	var minX, minY, maxX, maxY int
	c := func(x, y int) bool {
		if first {
			minX = x
			minY = y
			first = false
		}
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		return true
	}
	line(l.From.X, l.From.Y, l.To.X, l.To.Y, c)

	return image.Rect(0, 0, maxX-minX, maxY-minY)
}
