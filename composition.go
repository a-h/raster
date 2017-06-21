package raster

import (
	"image"
	"image/draw"

	"github.com/a-h/raster/affine"
)

// Composable represents a shape which can be combined with other shapes.
// Circle, Line, Polygon, Square and Text all implement this interface.
type Composable interface {
	Draw(img draw.Image)
	Bounds() image.Rectangle
}

type Composition struct {
	Position       image.Point
	Components     []Composable
	cache          *image.RGBA
	Transformation affine.Transformation
}

func NewComposition(position image.Point, components ...Composable) *Composition {
	return &Composition{
		Position:       position,
		Components:     components,
		Transformation: affine.NewTransformation(affine.IdentityMatrix),
	}
}

func (c *Composition) Draw(img draw.Image) {
	// Draw on a temporary canvas.
	// Cache the base image.
	if c.cache == nil {
		c.cache = image.NewRGBA(c.Bounds())
		for _, component := range c.Components {
			component.Draw(c.cache)
		}
	}

	// Apply the composition's transformations each time.
	t := affine.NewTranslationTransformation(c.Position.X, c.Position.Y)
	t = t.Combine(c.Transformation)
	for y := 0; y < c.cache.Bounds().Dy(); y++ {
		for x := 0; x < c.cache.Bounds().Dx(); x++ {
			transformedPoint := t.Apply(image.Point{x, y})

			img.Set(transformedPoint.X, transformedPoint.Y, c.cache.At(x, y))
		}
	}
}

func (c *Composition) Bounds() image.Rectangle {
	//TODO: Test the effect of the affine transformations.
	maxX := 0
	maxY := 0

	for _, component := range c.Components {
		if component.Bounds().Dx() > maxX {
			maxX = component.Bounds().Dx()
		}
		if component.Bounds().Dy() > maxY {
			maxY = component.Bounds().Dy()
		}
	}

	return image.Rect(0, 0, maxX+1, maxY+1)
}
