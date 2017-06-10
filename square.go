package raster

import (
	"image"
	"image/color"
)

type Square struct {
	Position image.Point
	Size     int
}

func NewSquare(position image.Point, size int) Square {
	return Square{
		Position: position,
		Size:     size,
	}
}

func (s Square) Draw(img *image.RGBA, o color.RGBA) (outline []image.Point) {
	a := image.Point{s.Position.X, s.Position.Y}
	b := image.Point{s.Position.X + s.Size, s.Position.Y}
	c := image.Point{s.Position.X + s.Size, s.Position.Y + s.Size}
	d := image.Point{s.Position.X, s.Position.Y + s.Size}

	top := NewLine(a.X, a.Y, b.X, b.Y)
	right := NewLine(b.X, b.Y, c.X, c.Y)
	bottom := NewLine(c.X, c.Y, d.X, d.Y)
	left := NewLine(d.X, d.Y, a.X, a.Y)

	points := top.Draw(img, o)
	points = concat(points, right.Draw(img, o))
	points = concat(points, bottom.Draw(img, o))
	return concat(points, left.Draw(img, o))
}

func concat(a []image.Point, b []image.Point) []image.Point {
	c := a
	for _, p := range b {
		c = append(c, p)
	}
	return c
}
