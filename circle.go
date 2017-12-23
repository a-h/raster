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

// Draw draws the element to the img, img could be an image.RGBA* or screen buffer.
func (c Circle) Draw(img draw.Image) image.Rectangle {
	bounds := image.Rect(c.Center.X-c.Radius-2, c.Center.Y-c.Radius-2, c.Center.X+c.Radius+2, c.Center.Y+c.Radius+2)
	for iy := bounds.Min.Y; iy < bounds.Max.Y; iy++ {
		// Work out from the left.
		foundBorder := false
		for ix := bounds.Min.X; ix < bounds.Max.X-c.Radius; ix++ {
			width := c.Center.X - ix
			height := c.Center.Y - iy

			distanceFromCenter := math.Sqrt(float64(((width * width) + (height * height))))
			onRadius := int(distanceFromCenter) == c.Radius

			if onRadius {
				img.Set(ix, iy, c.OutlineColor)
				foundBorder = true
			}

			// We've gone past the radius.
			if !onRadius && foundBorder {
				break
			}
		}
		// Work in from the right.
		foundBorder = false
		for ix := bounds.Max.X; ix > bounds.Max.X-c.Radius; ix-- {
			width := c.Center.X - ix
			height := c.Center.Y - iy

			distanceFromCenter := math.Sqrt(float64(((width * width) + (height * height))))
			onRadius := int(distanceFromCenter) == c.Radius

			if onRadius {
				img.Set(ix, iy, c.OutlineColor)
				foundBorder = true
			}

			// We've gone past the radius.
			if !onRadius && foundBorder {
				break
			}
		}
	}
	return bounds
}

// Bounds is the size of the object.
func (c Circle) Bounds() image.Rectangle {
	diameter := (c.Radius * 2)
	return image.Rect(0, 0, diameter, diameter)
}
