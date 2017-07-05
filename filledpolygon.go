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

	return image.Rect(0, 0, maxX-minX, maxY-minY)
}

// Draw draws the filled polygon onto the image.
func (p FilledPolygon) Draw(img draw.Image) {
	// First draw into a subimage.
	subImage := image.Rectangle{
		Min: image.Point{img.Bounds().Dx(), img.Bounds().Dy()},
		Max: image.Point{},
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
	subpolygon := NewPolygon(p.OutlineColor, translatedPoints...)

	subpolygonBounds := subpolygon.Bounds()
	subpolygonHeight := subpolygonBounds.Dy()
	subpolygonWidth := subpolygonBounds.Dx()

	// Scan across.
	for y := 0; y <= subpolygonHeight; y++ {
		insidePolygon := false
		passedLines := map[*Line]interface{}{}
		for x := 0; x <= subpolygonWidth; x++ {
			// Fill the polygon.
			if insidePolygon {
				canvas.Set(x, y, p.FillColor)
			}
			for _, line := range subpolygon.Lines {
				// Skip lines we've already passed
				if _, passed := passedLines[line]; passed {
					continue
				}

				lineCrossesY := (line.From.Y < y && line.To.Y >= y) ||
					(line.To.Y < y && line.From.Y >= y)
				lineCrossesX := (line.From.X <= x && line.To.X >= x) ||
					(line.To.X <= x && line.From.X >= x)

				// Only bother looking up if the line crosses the and x y axis.
				if lineCrossesY && lineCrossesX {
					if line.ContainsPoint(image.Point{x, y}) {
						// Mark that we've crossed a boundary
						insidePolygon = !insidePolygon
						passedLines[line] = true
					}
				}
			}
		}
	}

	// Draw the borders.
	subpolygon.Draw(canvas)

	// Copy everything that isn't transparent from the canvas to the target image at the subImage position.
	draw.Draw(img, subImage, canvas, image.Point{}, draw.Over)
}
