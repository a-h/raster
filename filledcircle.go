package raster

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// FilledCircle represents a circle, defined by a radius.
type FilledCircle struct {
	Circle
	FillColor color.RGBA
}

// NewFilledCircle creates a new circle, with the specified radius, filled with the fillcolor.
func NewFilledCircle(center image.Point, radius int, outlineColor color.RGBA, fillColor color.RGBA) FilledCircle {
	fc := FilledCircle{
		FillColor: fillColor,
	}
	fc.Center = center
	fc.Radius = radius
	fc.OutlineColor = outlineColor
	return fc
}

// Draw draws the element to the img, img could be an image.RGBA* or screen buffer.
func (c FilledCircle) Draw(img draw.Image) {
	bounds := image.Rect(c.Center.X-c.Radius-2, c.Center.Y-c.Radius-2, c.Center.X+c.Radius+2, c.Center.Y+c.Radius+2)
	for ix := bounds.Min.X; ix < bounds.Max.X; ix++ {
		for iy := bounds.Min.Y; iy < bounds.Max.Y; iy++ {
			width := c.Center.X - ix
			height := c.Center.Y - iy

			distanceFromCenter := math.Sqrt(float64(((width * width) + (height * height))))
			if int(distanceFromCenter) == c.Radius {
				img.Set(ix, iy, c.OutlineColor)
			}
			if int(distanceFromCenter) < c.Radius {
				img.Set(ix, iy, c.FillColor)
			}
		}
	}
}
