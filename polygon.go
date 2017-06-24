package raster

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/a-h/raster/biggest"
	"github.com/a-h/raster/smallest"
)

// Polygon defines a shape made from multiple lines.
type Polygon struct {
	Vertices     []image.Point
	Lines        []*Line
	OutlineColor color.RGBA
}

// NewPolygon creates a polygon made from lines which meet at the provided points (vertices).
func NewPolygon(outlineColor color.RGBA, vertices ...image.Point) Polygon {
	lines := []*Line{}

	// Calculate the lines.
	previousVertex := vertices[0]
	for _, p := range vertices[1:] {
		lines = append(lines, NewLine(previousVertex, p, outlineColor))
		previousVertex = p
	}
	lines = append(lines, NewLine(previousVertex, vertices[0], outlineColor))

	return Polygon{
		Vertices:     []image.Point(vertices),
		Lines:        lines,
		OutlineColor: outlineColor,
	}
}

// Bounds returns the size of the polygon.
func (p Polygon) Bounds() image.Rectangle {
	minX := p.Lines[0].From.X
	maxX := p.Lines[0].To.X
	minY := p.Lines[0].From.Y
	maxY := p.Lines[0].To.Y

	for _, l := range p.Lines {
		minX = smallest.IntegerIn(minX, l.From.X, l.To.X)
		minY = smallest.IntegerIn(minY, l.From.Y, l.To.Y)
		maxX = biggest.IntegerIn(maxX, l.From.X, l.To.X)
		maxY = biggest.IntegerIn(maxY, l.From.Y, l.To.Y)
	}

	return image.Rect(0, 0, maxX-minX, maxY-minY)
}

// Draw draws the polygon to the image.
func (p Polygon) Draw(img draw.Image) {
	for _, l := range p.Lines {
		l.Draw(img)
	}
}
