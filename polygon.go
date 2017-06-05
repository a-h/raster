package raster

import (
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/a-h/linear/tolerance"
)

type Polygon struct {
	Vertices []image.Point
	Lines    []*Line
}

func NewPolygon(vertices ...image.Point) Polygon {
	lines := []*Line{}

	// Calculate the lines.
	previousVertex := vertices[0]
	for _, p := range vertices[1:] {
		lines = append(lines, NewLine(previousVertex.X, previousVertex.Y, p.X, p.Y))
		previousVertex = p
	}
	lines = append(lines, NewLine(previousVertex.X, previousVertex.Y, vertices[0].X, vertices[0].Y))

	return Polygon{
		Vertices: []image.Point(vertices),
		Lines:    lines,
	}
}

// Bounds returns the size of the polygon.
func (p Polygon) Bounds() image.Rectangle {
	minX := p.Lines[0].From.X
	maxX := p.Lines[0].To.X
	minY := p.Lines[0].From.Y
	maxY := p.Lines[0].To.Y

	for _, l := range p.Lines {
		minX = smallest(minX, l.From.X, l.To.X)
		minY = smallest(minY, l.From.Y, l.To.Y)
		maxX = biggest(maxX, l.From.X, l.To.X)
		maxY = biggest(maxY, l.From.Y, l.To.Y)
	}

	return image.Rect(0, 0, maxX-minX, maxY-minY)
}

func smallest(ints ...int) int {
	sort.Ints(ints)
	return ints[0]
}

func biggest(ints ...int) int {
	sort.Ints(ints)
	return ints[len(ints)-1]
}

func (p Polygon) Points() (points []image.Point) {
	for _, l := range p.Lines {
		for _, p := range l.Points() {
			points = append(points, p)
		}
	}
	return
}

func (p Polygon) Draw(img *image.RGBA, o color.RGBA) []image.Point {
	points := p.Points()
	for _, p := range points {
		img.Set(p.X, p.Y, o)
	}
	return points
}

func (p Polygon) DrawFill(img *image.RGBA, o color.RGBA, f color.RGBA) []image.Point {
	// First draw into a subimage.
	subImage := image.Rectangle{
		Min: p.Vertices[0],
		Max: p.Vertices[0],
	}

	for _, pt := range p.Vertices {
		if pt.X < subImage.Min.X {
			subImage.Min.X = pt.X
		}
		if pt.Y < subImage.Min.Y {
			subImage.Min.Y = pt.Y
		}
		if pt.X > subImage.Max.X {
			subImage.Max.X = pt.X
		}
		if pt.Y > subImage.Max.Y {
			subImage.Max.Y = pt.Y
		}
	}

	// We have the origin of the image at topLeft, and know its size.
	// We can now print the image onto a canvas with a blank background.
	offsetX := subImage.Min.X
	offsetY := subImage.Min.Y
	canvas := image.NewRGBA(image.Rect(0, 0, (subImage.Max.X-offsetX)+1, (subImage.Max.Y-offsetY)+1))

	// Translate the points from the global to local coordinate space.
	translatedPoints := make([]image.Point, len(p.Vertices))
	for i, p := range p.Vertices {
		translatedPoints[i] = image.Point{X: p.X - offsetX, Y: p.Y - offsetY}
	}

	// Create the subpolygon.
	subpolygon := NewPolygon(translatedPoints...)

	subpolygonBounds := subpolygon.Bounds()
	subpolygonHeight := subpolygonBounds.Dy()
	subpolygonWidth := subpolygonBounds.Dx()

	// Sorted lines
	sortedLines := subpolygon.Lines
	/*
		sort.Slice(sortedLines, func(i, j int) bool {
			if sortedLines[i].From.X != sortedLines[j].From.X {
				return sortedLines[i].From.X < sortedLines[j].From.X
			}
			return sortedLines[i].From.Y < sortedLines[j].From.Y
		})
	*/

	// Scan across.
	for y := 0; y <= subpolygonHeight; y++ {
		insidePolygon := false
		passedLines := map[*Line]interface{}{}
		for x := 0; x <= subpolygonWidth; x++ {
			for _, line := range sortedLines {
				// Skip lines we've already passed
				if _, passed := passedLines[line]; passed {
					continue
				}
				// If the line crosses the y axis.
				if (line.From.Y < y && line.To.Y >= y) ||
					(line.To.Y < y && line.From.Y >= y) {
					if isPointOnLine2(line, image.Point{x, y}) {
						insidePolygon = !insidePolygon
						passedLines[line] = true
					}
				}
			}
			// Fill the polygon.
			if insidePolygon {
				canvas.Set(x, y, f)
			}
		}
	}

	// Draw the borders.
	subpolygon.Draw(canvas, o)

	// Copy everything that isn't transparent from the canvas to the target image at the subImage position.
	return drawNonTransparent(img, subImage, canvas, image.Point{})
}

func isPointOnLine2(l *Line, c image.Point) bool {
	distFromAToB := distance(l.From, l.To)

	distFromAToC := distance(l.From, c)
	distFromBToC := distance(l.To, c)
	distanceIncludingC := distFromAToC + distFromBToC

	return tolerance.IsWithin(distFromAToB, distanceIncludingC, 0.1)
}

func distance(p1, p2 image.Point) float64 {
	distY := math.Abs(float64(p2.Y - p1.Y))
	if p1.X == p2.X {
		// Vertical
		return distY
	}

	distX := math.Abs(float64(p2.X - p1.X))
	if p1.Y == p2.Y {
		// Horizontal
		return distX
	}

	a2 := distX * distX
	b2 := distY * distY
	return math.Sqrt(float64(a2) + float64(b2))
}

func isPointOnLine(l *Line, c image.Point) bool {
	// if AC is horizontal
	if l.From.X == l.To.X {
		return l.To.X == c.X
	}
	// if AC is vertical.
	if l.From.Y == l.To.Y {
		return l.To.Y == c.Y
	}
	// match the gradients
	return (l.From.X-c.X)*(l.From.Y-c.Y) == (c.X-l.To.X)*(c.Y-l.To.Y)
}

func drawNonTransparent(dst *image.RGBA, r image.Rectangle, src *image.RGBA, sp image.Point) []image.Point {
	points := []image.Point{}
	for srcX := sp.X; srcX < src.Bounds().Dx(); srcX++ {
		for srcY := sp.Y; srcY < src.Bounds().Dy(); srcY++ {
			dstX := r.Min.X + srcX
			dstY := r.Min.Y + srcY

			dstColor := src.At(srcX, srcY)
			if !isTransparent(dstColor) {
				points = append(points, image.Point{dstX, dstY})
				dst.Set(dstX, dstY, dstColor)
			}
		}
	}
	return points
}

func isTransparent(c color.Color) bool {
	r1, g1, b1, a1 := c.RGBA()
	r2, g2, b2, a2 := color.Transparent.RGBA()
	if r1 != r2 {
		return false
	}
	if g1 != g2 {
		return false
	}
	if b1 != b2 {
		return false
	}
	if a1 != a2 {
		return false
	}
	return true
}

func fillPoints(img *image.RGBA, polygon Polygon) []image.Point {
	points := []image.Point{}
	for y := 0; y < img.Bounds().Dy(); y++ {
		scan := ScanLine(y, polygon)
		for x, intersections := range scan {
			if intersections%2 > 0 {
				// We're inside the polygon, because we've intersected an odd number of times.
				points = append(points, image.Point{x, y})
			}
		}
	}
	return points
}

// IsEdge returns true when two lines are next to each other in the Polygon list.
func (p Polygon) IsEdge(a image.Point) (edge bool, linesWhichMeet []*Line) {
	for _, l := range p.Lines {
		if l.ContainsPoint(a) {
			linesWhichMeet = append(linesWhichMeet, l)
		}
	}
	edge = len(linesWhichMeet) > 0
	return
}
