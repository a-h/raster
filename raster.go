package raster

import (
	"image"
	"image/color"
	"math"
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
func DrawLine(img *image.RGBA, fromX, fromY int, toX, toY int, c color.RGBA) {
	// We're moving from fromX to toX, so make sure they're in the right order.
	if toX < fromX {
		toX, toY, fromX, fromY = fromX, fromY, toX, toY
	}

	// Vertical line.
	if fromX == toX {
		for y := fromY; y <= toY; y++ {
			img.Set(fromX, y, c)
		}
		return
	}

	// Horizontal line.
	if fromY == toY {
		for x := fromX; x <= toX; x++ {
			img.Set(x, fromY, c)
		}
		return
	}

	// It's a slope.
	rise := toY - fromY
	run := toX - fromX
	m := float64(rise) / float64(run)

	y := float64(fromY)
	for x := fromX; x <= toX; x++ {
		img.Set(x, int(y), c)
		y += m
	}
}
