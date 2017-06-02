package raster

import (
	"image"
	"image/color"
	"math"
	"sort"
)

// DrawCircle draws a Circle onto the image at the x, y coordinates.
func DrawCircle(img *image.RGBA, x, y int, radius int, c color.RGBA) {
	bounds := image.Rect(x-radius-2, y-radius-2, x+radius+2, y+radius+2)
	for ix := bounds.Min.X; ix < bounds.Max.X; ix++ {
		for iy := bounds.Min.Y; iy < bounds.Max.Y; iy++ {
			width := x - ix
			height := y - iy

			distanceFromCenter := math.Sqrt(float64(((width * width) + (height * height))))
			if int(distanceFromCenter) == radius {
				img.Set(ix, iy, c)
			}
		}
	}
}

// DrawDisc draws a filled circle onto the image at the x, y coordinates.
func DrawDisc(img *image.RGBA, x, y int, radius int, c color.RGBA) {
	bounds := image.Rect(x-radius-2, y-radius-2, x+radius+2, y+radius+2)
	for ix := bounds.Min.X; ix < bounds.Max.X; ix++ {
		for iy := bounds.Min.Y; iy < bounds.Max.Y; iy++ {
			width := x - ix
			height := y - iy

			distanceFromCenter := math.Sqrt(float64(((width * width) + (height * height))))
			if int(distanceFromCenter) <= radius {
				img.Set(ix, iy, c)
			}
		}
	}
}

// DrawLine draws a line circle onto the image starting at the fromX and fromY coordinates to the
// toX, toY coordinates.
func DrawLine(img *image.RGBA, fromX, fromY int, toX, toY int, c color.RGBA) []image.Point {
	points := Line(fromX, fromY, toX, toY)
	for _, p := range points {
		img.Set(p.X, p.Y, c)
	}
	return points
}

func Line(fromX, fromY int, toX, toY int) (points []image.Point) {
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

	if rise > run {
		m := float64(run) / float64(rise)
		x := float64(fromX)
		for y := fromY; y < toY; y++ {
			points = append(points, image.Point{int(x), y})
			x += m
		}
	} else {
		m := float64(rise) / float64(run)
		y := float64(fromY)
		for x := fromX; x < toX; x++ {
			points = append(points, image.Point{x, int(y)})
			y += m
		}
	}
	points = append(points, image.Point{toX, toY})
	return
}

func DrawPolygon(img *image.RGBA, c color.RGBA, vertices ...image.Point) LineMap {
	lm := Polygon(vertices...)
	for _, p := range lm.Points() {
		img.Set(p.X, p.Y, c)
	}
	return lm
}

func Polygon(vertices ...image.Point) LineMap {
	lm := NewLineMap()
	previousVertex := vertices[0]
	for _, p := range vertices[1:] {
		lm.AddLine(Line(previousVertex.X, previousVertex.Y, p.X, p.Y))
		previousVertex = p
	}
	lm.AddLine(Line(previousVertex.X, previousVertex.Y, vertices[0].X, vertices[0].Y))
	return lm
}

func DrawFilledPolygon(img *image.RGBA, outline color.RGBA, fill color.RGBA, points ...image.Point) {
	// Find the bounding box of the target area.
	subImage := image.Rectangle{
		Min: points[0],
		Max: points[0],
	}

	for _, p := range points {
		if p.X < subImage.Min.X {
			subImage.Min.X = p.X
		}
		if p.Y < subImage.Min.Y {
			subImage.Min.Y = p.Y
		}
		if p.X > subImage.Max.X {
			subImage.Max.X = p.X
		}
		if p.Y > subImage.Max.Y {
			subImage.Max.Y = p.Y
		}
	}

	// We have the origin of the image at topLeft, and know its size.
	// We can now print the image onto a canvas with a blank background.
	offsetX := subImage.Min.X
	offsetY := subImage.Min.Y
	canvas := image.NewRGBA(image.Rect(0, 0, (subImage.Max.X-offsetX)+1, (subImage.Max.Y-offsetY)+1))

	// Translate the points into the local space.
	translatedPoints := make([]image.Point, len(points))
	for i, p := range points {
		translatedPoints[i] = image.Point{X: p.X - offsetX, Y: p.Y - offsetY}
	}

	outlinePoints := DrawPolygon(canvas, outline, translatedPoints...)

	// Use a simple ray algorithm to fill.
	FillBetweenLines(canvas, fill, outlinePoints)

	// draw.Draw(img, subImage, canvas, image.Point{}, draw.Over)
	// draw.DrawMask(img, subImage, canvas, image.Point{}, &image.Uniform{color.Transparent}, subImage.Min, draw.Over)
	DrawNonTransparent(img, subImage, canvas, image.Point{})
}

// Raycast counts the intersections.
func Raycast(current image.Point, r image.Rectangle, outline LineMap) int {
	// Work out the shortest direction.
	distanceToTop := current.Y
	distanceToBottom := r.Dy() - current.Y
	distanceToLeft := current.X
	distanceToRight := r.Dx() - current.X

	values := []int{distanceToTop, distanceToBottom, distanceToLeft, distanceToRight}
	sort.Ints(values)
	smallest := values[0]

	if distanceToTop == smallest {
		return RaycastUp(current, r, outline)
	}
	if distanceToBottom == smallest {
		return RaycastDown(current, r, outline)
	}
	if distanceToLeft == smallest {
		return RaycastLeft(current, r, outline)
	}
	return RaycastRight(current, r, outline)
}

func RaycastLeft(current image.Point, r image.Rectangle, outline LineMap) int {
	count := 0
	intersected := make(map[int]interface{})
	for x := current.X; x >= 0; x-- {
		if intersections, ok := outline.Lookup[image.Point{x, current.Y}]; ok {
			for _, intersection := range intersections {
				// If we've not intersected before, then it counts.
				if _, ok := intersected[intersection]; !ok {
					count++
					intersected[intersection] = true
				}
			}
		}
	}
	return count
}

func RaycastRight(current image.Point, r image.Rectangle, outline LineMap) int {
	count := 0
	intersected := make(map[int]interface{})
	for x := current.X; x < r.Dx(); x++ {
		if intersections, ok := outline.Lookup[image.Point{x, current.Y}]; ok {
			for _, intersection := range intersections {
				// If we've not intersected before, then it counts.
				if _, ok := intersected[intersection]; !ok {
					count++
					intersected[intersection] = true
				}
			}
		}
	}
	return count
}

func RaycastUp(current image.Point, r image.Rectangle, outline LineMap) int {
	count := 0
	intersected := make(map[int]interface{})
	for y := current.Y; y >= 0; y-- {
		// Get the line
		if intersections, ok := outline.Lookup[image.Point{current.X, y}]; ok {
			for _, intersection := range intersections {
				// If we've not intersected before, then it counts.
				if _, ok := intersected[intersection]; !ok {
					count++
					intersected[intersection] = true
				}
			}
		}
	}
	return count
}

func RaycastDown(current image.Point, r image.Rectangle, outline LineMap) int {
	count := 0
	intersected := make(map[int]interface{})
	for y := current.Y; y < r.Dy(); y++ {
		if intersections, ok := outline.Lookup[image.Point{current.X, y}]; ok {
			for _, intersection := range intersections {
				// If we've not intersected before, then it counts.
				if _, ok := intersected[intersection]; !ok {
					count++
					intersected[intersection] = true
				}
			}
		}
	}
	return count
}

func FillBetweenLines(img *image.RGBA, c color.Color, outline LineMap) {
	// Use a Ray Casting algorithm.
	// https://en.wikipedia.org/wiki/Point_in_polygon
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			p := image.Point{x, y}
			_, onEdge := outline.Lookup[p]
			if onEdge {
				continue
			}
			intersections := Raycast(p, img.Bounds(), outline)

			if intersections > 0 && intersections%2 != 0 {
				// We're inside the polygon.
				img.Set(x, y, c)
			}
		}
	}
}

func arraysAreEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	bmap := make(map[int]interface{})
	for _, b1 := range b {
		bmap[b1] = struct{}{}
	}
	for _, a1 := range a {
		_, ok := bmap[a1]
		if !ok {
			return false
		}
	}
	return true
}

func DrawNonTransparent(dst *image.RGBA, r image.Rectangle, src *image.RGBA, sp image.Point) {
	for srcX := sp.X; srcX < src.Bounds().Dx(); srcX++ {
		for srcY := sp.Y; srcY < src.Bounds().Dy(); srcY++ {
			dstX := r.Min.X + srcX
			dstY := r.Min.Y + srcY

			srcColor := src.At(srcX, srcY)
			if !isTransparent(srcColor) {
				dst.Set(dstX, dstY, srcColor)
			}
		}
	}
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

func NewLineMap() LineMap {
	return LineMap{
		Lines:  0,
		Lookup: make(map[image.Point][]int),
	}
}

// LineMap tracks which pixels belong to particular lines.
type LineMap struct {
	// Lines is the numbmer of lines in the map.
	Lines int
	// Find which line or lines a particular point is part of.
	Lookup map[image.Point][]int
}

// AddLine adds all the points of a line and returns the index of the line, e.g.
// you add the first set of points which make up the top of a square and get back
// 1. This is not thread-safe.
func (lm *LineMap) AddLine(points []image.Point) {
	lm.Lines++
	currentLine := lm.Lines
	for _, p := range points {
		a, _ := lm.Lookup[p]
		lm.Lookup[p] = append(a, currentLine)
	}
}

// Points provides a sorted slice of all points in the map.
func (lm *LineMap) Points() []image.Point {
	// Sort the points.
	points := make([]image.Point, len(lm.Lookup))
	for p := range lm.Lookup {
		points = append(points, p)
	}
	sort.Slice(points, func(i, j int) bool {
		if points[i].X != points[j].X {
			return points[i].X < points[j].X
		}
		return points[i].Y < points[j].Y
	})
	return points
}