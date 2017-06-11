package raster

import (
	"image"
	"image/color"
)

type Square struct {
	Position     image.Point
	Size         int
	OutlineColor color.RGBA
}

func NewSquare(x, y int, size int, outlineColor color.RGBA) Square {
	return Square{
		Position:     image.Point{x, y},
		Size:         size,
		OutlineColor: outlineColor,
	}
}

func (s Square) Draw(img *image.RGBA) (outline []image.Point) {
	a := image.Point{s.Position.X, s.Position.Y}
	b := image.Point{s.Position.X + s.Size, s.Position.Y}
	c := image.Point{s.Position.X + s.Size, s.Position.Y + s.Size}
	d := image.Point{s.Position.X, s.Position.Y + s.Size}

	top := NewLine(a.X, a.Y, b.X, b.Y, s.OutlineColor)
	right := NewLine(b.X, b.Y, c.X, c.Y, s.OutlineColor)
	bottom := NewLine(c.X, c.Y, d.X, d.Y, s.OutlineColor)
	left := NewLine(d.X, d.Y, a.X, a.Y, s.OutlineColor)

	points := top.Draw(img)
	points = append(points, right.Draw(img)...)
	points = append(points, bottom.Draw(img)...)
	return append(points, left.Draw(img)...)
}

// Bounds returns the size of the object.
func (s Square) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.Size, s.Size)
}
