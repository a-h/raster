package sparse

import (
	"image"
	"image/color"
)

// Image holds the written pixels in a map instead of creating memory to hold all pixels.
type Image struct {
	bounds image.Rectangle
	Drawn  map[image.Point]color.Color
}

// NewImage creates a sparse image.
func NewImage(bounds image.Rectangle) *Image {
	return &Image{
		bounds: bounds,
		Drawn:  make(map[image.Point]color.Color),
	}
}

// ColorModel returns the Image's color model.
func (img *Image) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (img *Image) Bounds() image.Rectangle {
	return img.bounds
}

// DrawnBounds returns the bounds of the image which contain information, not
// the theoretical size of the image.
func (img *Image) DrawnBounds() image.Rectangle {
	if len(img.Drawn) == 0 {
		return image.Rect(0, 0, 0, 0)
	}

	b := img.Bounds()
	minX, minY := b.Dx(), b.Dy()
	maxX, maxY := 0, 0
	for p := range img.Drawn {
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
	}
	return image.Rect(minX, minY, maxX, maxY)
}

// Set the pixel with x, y coordinates to the color c.
func (img *Image) Set(x, y int, c color.Color) {
	img.Drawn[image.Point{x, y}] = c
}

// At returns the color of the pixel at (x, y).
func (img *Image) At(x, y int) color.Color {
	p := image.Point{x, y}
	if c, ok := img.Drawn[p]; ok {
		return c
	}
	return color.RGBA{}
}
