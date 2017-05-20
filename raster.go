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
	// Vertical line.
	if fromX == toX {
		if toY < fromY {
			toX, toY, fromX, fromY = fromX, fromY, toX, toY
		}
		for y := fromY; y <= toY; y++ {
			img.Set(fromX, y, c)
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

func DrawPolygon(img *image.RGBA, c color.RGBA, points ...image.Point) {
	previousPoint := points[0]
	for _, p := range points[1:] {
		DrawLine(img, previousPoint.X, previousPoint.Y, p.X, p.Y, c)
		previousPoint = p
	}
	DrawLine(img, previousPoint.X, previousPoint.Y, points[0].X, points[0].Y, c)
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
	canvas := image.NewRGBA(image.Rect(0, 0, subImage.Max.X-offsetX+1, subImage.Max.Y-offsetY+1))

	// Translate the points into the local space.
	translatedPoints := make([]image.Point, len(points))
	for i, p := range points {
		translatedPoints[i] = image.Point{X: p.X - offsetX, Y: p.Y - offsetY}
	}

	previousPoint := translatedPoints[0]
	for _, p := range translatedPoints[1:] {
		DrawLine(canvas, previousPoint.X, previousPoint.Y, p.X, p.Y, outline)
		previousPoint = p
	}
	DrawLine(canvas, previousPoint.X, previousPoint.Y, translatedPoints[0].X, translatedPoints[0].Y, outline)

	// Use a simple ray algorithm to fill.
	FillBetweenLines(canvas, fill, translatedPoints)

	// draw.Draw(img, subImage, canvas, image.Point{}, draw.Over)
	// draw.DrawMask(img, subImage, canvas, image.Point{}, &image.Uniform{color.Transparent}, subImage.Min, draw.Over)
	DrawNonTransparent(img, subImage, canvas, image.Point{})
}

func FillBetweenLines(img *image.RGBA, c color.Color, vertices []image.Point) {
	vertexMap := make(map[image.Point]interface{})
	for _, v := range vertices {
		vertexMap[v] = false
	}
	// Use a Ray Casting algorithm, it won't work properly until I make sure that vertices are ignored.
	// https://en.wikipedia.org/wiki/Point_in_polygon
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			if isTransparent(img.At(x, y)) {
				// Look to the right edge to see if we're inside the polygon.
				intersections := 0
				for ix := x; ix < img.Bounds().Dx(); ix++ {
					if !isTransparent(img.At(ix, y)) {
						_, isVertex := vertexMap[image.Point{ix, y}]
						if !isVertex {
							intersections++
						}
					}
				}

				if intersections%2 != 0 {
					// We're inside the polygon.
					img.Set(x, y, c)
				}
			}
		}
	}
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
