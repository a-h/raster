package raster

import (
	"image"
	"image/color"
	"image/draw"
)

// A Square has a position, size and outline color.
type Square struct {
	Position     image.Point
	Size         int
	OutlineColor color.RGBA
}

// NewSquare creates a new square. The position represents the top left coordinate.
func NewSquare(position image.Point, size int, outlineColor color.RGBA) Square {
	return Square{
		Position:     position,
		Size:         size,
		OutlineColor: outlineColor,
	}
}

func (s Square) Draw(img draw.Image) {
	a := image.Point{s.Position.X, s.Position.Y}
	b := image.Point{s.Position.X + s.Size, s.Position.Y}
	c := image.Point{s.Position.X + s.Size, s.Position.Y + s.Size}
	d := image.Point{s.Position.X, s.Position.Y + s.Size}

	top := NewLine(a, b, s.OutlineColor)
	right := NewLine(b, c, s.OutlineColor)
	bottom := NewLine(c, d, s.OutlineColor)
	left := NewLine(d, a, s.OutlineColor)

	top.Draw(img)
	right.Draw(img)
	bottom.Draw(img)
	left.Draw(img)
}

// Bounds returns the size of the object.
func (s Square) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.Size, s.Size)
}
