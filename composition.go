package raster

import "image"

type Composable interface {
	Draw(img *image.RGBA) []image.Point
	Bounds() image.Rectangle
}

type Composition struct {
	Position   image.Point
	Components []Composable
}

func NewComposition(position image.Point, components ...Composable) *Composition {
	return &Composition{
		Position:   position,
		Components: components,
	}
}

func (c *Composition) Draw(img *image.RGBA) {
	// Draw on a temporary canvas.
	canvas := image.NewRGBA(c.Bounds())
	for _, component := range c.Components {
		component.Draw(canvas)
	}

	// Moves the object to the position.
	translationMatrix := []float64{
		1, 0, float64(c.Position.X),
		0, 1, float64(c.Position.Y),
	}

	// Does nothing to the object.
	/*
		identityMatrix := []float64{
			1, 0, 0,
			0, 1, 0,
		}
	*/

	/*
		rotateBy := float64(30)
		rotationMatrix := []float64{
			math.Cos(rotateBy), -math.Sin(rotateBy), 0,
			math.Sin(rotateBy), math.Cos(rotateBy), 0,
		}*/
	t := NewTransformation(translationMatrix)
	for y := 0; y < canvas.Bounds().Dy(); y++ {
		for x := 0; x < canvas.Bounds().Dx(); x++ {
			transformedPoint := t.Apply(image.Point{x, y})

			img.Set(transformedPoint.X, transformedPoint.Y, canvas.At(x, y))
		}
	}
}

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

type Transformation struct {
	a, b, c float64
	p, q, r float64
	u, v, w float64
}

func NewTransformation(matrix []float64) Transformation {
	return Transformation{
		a: matrix[0],
		b: matrix[1],
		c: matrix[2],
		p: matrix[3],
		q: matrix[4],
		r: matrix[5],
		u: 0,
		v: 0,
		w: 1,
	}
}

func (t Transformation) Apply(point image.Point) image.Point {
	// See https://en.wikipedia.org/wiki/Matrix_multiplication#Matrix_product_.28two_matrices.29
	// Square matrix and column vector (the point)
	x := float64(point.X)
	y := float64(point.Y)
	z := float64(1)

	x1 := (t.a * x) + (t.b * y) + (t.c * z)
	y1 := (t.p * x) + (t.q * y) + (t.r * z)

	return image.Point{int(x1), int(y1)}
}
