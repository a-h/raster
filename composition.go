package raster

import (
	"image"
	"image/draw"

	"github.com/a-h/raster/biggest"
	"github.com/a-h/raster/smallest"

	"github.com/a-h/raster/affine"
	"github.com/a-h/raster/sparse"
)

// Composable represents a shape which can be combined with other shapes.
// Circle, Line, Polygon, Square and Text all implement this interface.
type Composable interface {
	// Draw draws the element to the img, img could be an image.RGBA* or screen buffer.
	Draw(img draw.Image) image.Rectangle
	Bounds() image.Rectangle
}

// Composition returns the position and components which make it up, and a transformation
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
// It returns the area actually drawn out on the image.
func (c *Composition) Draw(img draw.Image) image.Rectangle {
	// Draw on a temporary canvas.
	// Cache the base image.
	if c.cache == nil {
		c.cache = sparse.NewImage(c.Bounds())
		for _, component := range c.Components {
			component.Draw(c.cache)
		}
	}

	// Apply the composition's transformations each time.
	minX, minY, maxX, maxY := c.Position.X, c.Position.Y, 0, 0
	for position, color := range c.cache.Drawn {
		transformedPoint := c.Transformation.Apply(position)

		x, y := transformedPoint.X+c.Position.X, transformedPoint.Y+c.Position.Y
		minX = smallest.IntegerIn(minX, x)
		maxX = biggest.IntegerIn(maxX, x)
		minY = smallest.IntegerIn(minY, y)
		maxY = biggest.IntegerIn(maxY, y)

		img.Set(x, y, color)
	}

	return image.Rect(minX, minY, maxX+1, maxY+1)
}

// Bounds provides the area of the composition prior to affine transformations being
// applied.
func (c *Composition) Bounds() image.Rectangle {
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
