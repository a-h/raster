package raster

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// A Text element contains a Position, Text to write and foreground text color.
type Text struct {
	Position image.Point
	Text     string
	Color    color.RGBA
}

// NewText creates a text element at the specified position.
func NewText(position image.Point, text string, c color.RGBA) Text {
	return Text{
		Position: position,
		Text:     text,
		Color:    c,
	}
}

// Draw draws the element to the img, img could be an image.RGBA* or screen buffer.
func (t Text) Draw(img draw.Image) image.Rectangle {
	// See https://stackoverflow.com/questions/38299930/how-to-add-a-simple-text-label-to-an-image-in-go
	point := fixed.Point26_6{
		X: fixed.I(t.Position.X),
		Y: fixed.I(t.Position.Y + 13), // Fonts are drawn from the base point, not the top left.
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(colornames.White),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(t.Text)

	b, _ := d.BoundString(t.Text)
	return image.Rect(t.Position.X+b.Min.X.Round(), t.Position.Y+b.Min.Y.Round(),
		t.Position.X+b.Max.X.Round(), t.Position.Y+b.Max.Y.Round())
}

// Bounds returns the size of the object.
func (t Text) Bounds() image.Rectangle {
	d := &font.Drawer{
		Src:  image.NewUniform(colornames.White),
		Face: basicfont.Face7x13,
	}

	b, _ := d.BoundString(t.Text)
	return image.Rect(b.Min.X.Round(), b.Min.Y.Round(), b.Max.X.Round(), b.Max.Y.Round())
}
