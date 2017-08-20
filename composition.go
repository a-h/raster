package raster

import (
	"image"
	"image/draw"

	"github.com/a-h/raster/affine"
	"github.com/a-h/raster/sparse"
)

// Composable represents a shape which can be combined with other shapes.
// Circle, Line, Polygon, Square and Text all implement this interface.
type Composable interface {
	// Draw draws the element to the img, img could be an image.RGBA* or screen buffer.
	Draw(img draw.Image)
	Bounds() image.Rectangle
}

// A Composition has a position, components whic make it up, and a transformation
// that can be applied to it to move, scale, or rotate all of the elements.
type Composition struct {
	Position       image.Point
	Components     []Composable
	cache          *sparse.Image
	Transformation affine.Transformation
}

// NewComposition creates a composition for rendering at the specific point. The components must
// be provided with reference to the coordinates of the composition's top left corner.
func NewComposition(position image.Point, components ...Composable) *Composition {
	return &Composition{
		Position:       position,
		Components:     components,
		Transformation: affine.NewTransformation(affine.IdentityMatrix),
	}
}

// Draw draws the element to the img, img could be an image.RGBA* or screen buffer.
func (c *Composition) Draw(img draw.Image) {
	// Draw on a temporary canvas.
	// Cache the base image.
	if c.cache == nil {
		c.cache = sparse.NewImage(c.Bounds())
		for _, component := range c.Components {
			component.Draw(c.cache)
		}
	}

	// Apply the composition's transformations each time.
	t := affine.NewTranslationTransformation(c.Position.X, c.Position.Y)
	t = t.Combine(c.Transformation)

	for position, color := range c.cache.Drawn {
		transformedPoint := t.Apply(position)
		img.Set(transformedPoint.X, transformedPoint.Y, color)
	}
}

// Bounds provides the area of the composition.
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
