package raster

import (
	"image"
	"image/color"
	"math"
)

// Circle represents a circle, defined by a radius.
type Circle struct {
	Center       image.Point
	Radius       int
	OutlineColor color.RGBA
}

// NewCircle creates a new circle, with the specified radius.
func NewCircle(x, y int, radius int, outlineColor color.RGBA) Circle {
	return Circle{
		Center:       image.Point{x, y},
		Radius:       radius,
		OutlineColor: outlineColor,
	}
}

// Points returns the points which make up the circle.
func (c Circle) Points() (outline []image.Point, interior []image.Point) {
	bounds := image.Rect(c.Center.X-c.Radius-2, c.Center.Y-c.Radius-2, c.Center.X+c.Radius+2, c.Center.Y+c.Radius+2)
	for ix := bounds.Min.X; ix < bounds.Max.X; ix++ {
		for iy := bounds.Min.Y; iy < bounds.Max.Y; iy++ {
			width := c.Center.X - ix
			height := c.Center.Y - iy

			distanceFromCenter := math.Sqrt(float64(((width * width) + (height * height))))
			if int(distanceFromCenter) == c.Radius {
				outline = append(outline, image.Point{ix, iy})
			}
			if int(distanceFromCenter) < c.Radius {
				interior = append(interior, image.Point{ix, iy})
			}
		}
	}
	return outline, interior
}

// Draw draws the circle to the screen.
func (c Circle) Draw(img *image.RGBA) (outline []image.Point) {
	outline, _ = c.Points()
	for _, p := range outline {
		img.Set(p.X, p.Y, c.OutlineColor)
	}
	return outline
}

// DrawFilled draws the filled circle to the screen.
func (c Circle) DrawFilled(img *image.RGBA, o color.RGBA, f color.RGBA) (outline []image.Point, interior []image.Point) {
	outline, interior = c.Points()
	for _, p := range outline {
		img.Set(p.X, p.Y, o)
	}
	for _, p := range interior {
		img.Set(p.X, p.Y, f)
	}
	return nil, outline
}
