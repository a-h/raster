package raster

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/a-h/raster/biggest"
	"github.com/a-h/raster/smallest"
)

// FilledPolygon defines a shape made from multiple lines.
type FilledPolygon struct {
	Polygon
	FillColor color.RGBA
}

// NewFilledPolygon creates a polygon made from lines which meet at the provided points (vertices).
func NewFilledPolygon(outlineColor color.RGBA, fillColor color.RGBA, vertices ...image.Point) FilledPolygon {
	lines := []*Line{}

	// Calculate the lines.
	previousVertex := vertices[0]
	for _, p := range vertices[1:] {
		lines = append(lines, NewLine(previousVertex, p, outlineColor))
		previousVertex = p
	}
	lines = append(lines, NewLine(previousVertex, vertices[0], outlineColor))

	fp := FilledPolygon{
		FillColor: fillColor,
	}
	fp.OutlineColor = outlineColor
	fp.Vertices = []image.Point(vertices)
	fp.Lines = lines
	return fp
}

// Bounds returns the size of the polygon.
func (p FilledPolygon) Bounds() image.Rectangle {
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

	return image.Rect(minX, minY, maxX, maxY)
}

// Draw draws the filled polygon onto the image.
func (p FilledPolygon) Draw(img draw.Image) {
	// Create the outline.
	subpolygon := NewPolygon(p.OutlineColor, p.Vertices...)

	subpolygonBounds := subpolygon.Bounds()
	subpolygonHeight := subpolygonBounds.Dy()
	subpolygonWidth := subpolygonBounds.Dx()

	offsetX, offsetY := p.Vertices[0].X, p.Vertices[0].Y
	for _, pt := range p.Vertices[1:] {
		if pt.X < offsetX {
			offsetX = pt.X
		}
		if pt.Y < offsetY {
			offsetY = pt.Y
		}
	}

	// Scan across.
	for y := offsetY; y <= subpolygonHeight+offsetY+1; y++ {
		insidePolygon := false
		passedLines := map[*Line]interface{}{}
		for x := offsetX - 1; x <= subpolygonWidth+offsetX+1; x++ {
			// Fill the polygon.
			if insidePolygon {
				img.Set(x, y, p.FillColor)
			}
			for _, line := range subpolygon.Lines {
				// Skip lines we've already passed
				if _, passed := passedLines[line]; passed {
					continue
				}

				lineCrossesY := (line.From.Y < y && line.To.Y >= y) ||
					(line.To.Y < y && line.From.Y >= y)

				if lineCrossesY && line.ContainsPoint(image.Point{x, y}) {
					// Mark that we've crossed a boundary
					insidePolygon = !insidePolygon
					passedLines[line] = true
				}
			}
		}
	}

	// Draw the lines.
	subpolygon.Draw(img)
}
