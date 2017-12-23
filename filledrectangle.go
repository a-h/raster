package raster

import (
	"image"
	"image/color"
	"image/draw"
)

// A FilledRectangle has a position, size and outline color.
type FilledRectangle struct {
	Position     image.Point
	Width        int
	Height       int
	OutlineColor color.RGBA
	FillColor    color.RGBA
}

// NewFilledRectangle creates a new filled rectangle. The position represents the top left coordinate.
func NewFilledRectangle(position image.Point, width, height int, outline, fill color.RGBA) FilledRectangle {
	return FilledRectangle{
		Position:     position,
		Width:        width,
		Height:       height,
		OutlineColor: outline,
		FillColor:    fill,
	}
}

// Draw draws the element to the img, img could be an image.RGBA* or screen buffer.
func (r FilledRectangle) Draw(img draw.Image) image.Rectangle {
	for y := r.Position.Y; y < r.Position.Y+r.Height; y++ {
		for x := r.Position.X; x < r.Position.X+r.Width; x++ {
			img.Set(x, y, r.FillColor)
		}
	}

	a := image.Point{r.Position.X, r.Position.Y}
	b := image.Point{r.Position.X + r.Width, r.Position.Y}
	c := image.Point{r.Position.X + r.Width, r.Position.Y + r.Height}
	d := image.Point{r.Position.X, r.Position.Y + r.Height}

	top := NewLine(a, b, r.OutlineColor)
	right := NewLine(b, c, r.OutlineColor)
	bottom := NewLine(c, d, r.OutlineColor)
	left := NewLine(d, a, r.OutlineColor)

	top.Draw(img)
	right.Draw(img)
	bottom.Draw(img)
	left.Draw(img)

	return image.Rect(a.X, a.Y, d.X, d.Y)
}

// Bounds returns the size of the object.
func (r FilledRectangle) Bounds() image.Rectangle {
	return image.Rect(0, 0, r.Width, r.Height)
}
