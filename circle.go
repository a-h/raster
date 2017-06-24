package raster

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Circle represents a circle, defined by a radius.
type Circle struct {
	Center       image.Point
	Radius       int
	OutlineColor color.RGBA
}

// NewCircle creates a new circle, with the specified radius.
func NewCircle(center image.Point, radius int, outlineColor color.RGBA) Circle {
	return Circle{
		Center:       center,
		Radius:       radius,
		OutlineColor: outlineColor,
	}
}

// Draw draws the element to the image. The image could be an image.RGBA* or screen buffer.
func (c Circle) Draw(img draw.Image) {
	bounds := image.Rect(c.Center.X-c.Radius-2, c.Center.Y-c.Radius-2, c.Center.X+c.Radius+2, c.Center.Y+c.Radius+2)
	for ix := bounds.Min.X; ix < bounds.Max.X; ix++ {
		for iy := bounds.Min.Y; iy < bounds.Max.Y; iy++ {
			width := c.Center.X - ix
			height := c.Center.Y - iy

			distanceFromCenter := math.Sqrt(float64(((width * width) + (height * height))))
			if int(distanceFromCenter) == c.Radius {
				img.Set(ix, iy, c.OutlineColor)
			}
		}
	}
}

// Bounds is the size of the object.
func (c Circle) Bounds() image.Rectangle {
	diameter := (c.Radius * 2)
	return image.Rect(0, 0, diameter, diameter)
}
