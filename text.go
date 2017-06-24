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

type Text struct {
	Position image.Point
	Text     string
	Color    color.RGBA
}

func NewText(x, y int, text string, c color.RGBA) Text {
	return Text{
		Position: image.Point{x, y},
		Text:     text,
		Color:    c,
	}
}

func (t Text) Draw(img draw.Image) {
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
