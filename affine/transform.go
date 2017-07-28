package affine

import (
	"image"
	"math"

	"github.com/a-h/round"
)

// Transformation represents a 3x3 matrix used to carry out an affine transform.
type Transformation struct {
	a, b, c float64
	p, q, r float64
	u, v, w float64
}

// NewTransformation creates a custom transformation, based on recieving an array of 6 floating points.
// E.g. If passed in the the identity matrix, the transform does nothing.
// NewTransformation([]float64{
//    1, 0, 1,
//    0, 1, 1,
// })
// The last 3 values are set to 0, 0, 1 regardless of input size.
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

// IdentityMatrix defines a matrix which has no effect.
var IdentityMatrix = []float64{
	1, 0, 0,
	0, 1, 0,
	0, 0, 1,
}

// NewScaleTransformation applies a scaling factor to the width and height, e.g.
// a width of 0.5 would be half the size. Numbers greater than one will be ignored
// as inputs.
func NewScaleTransformation(width, height float64) Transformation {
	if width > 1 || height > 1 {
		return NewTransformation(IdentityMatrix)
	}
	return NewTransformation([]float64{
		width, 0, 0,
		0, height, 0,
	})
}

// NewTranslationTransformation moves the point elsewhere.
func NewTranslationTransformation(x, y int) Transformation {
	return NewTransformation([]float64{
		1, 0, float64(x),
		0, 1, float64(y),
	})
}

// NewReflectionTransformation mirrors the image.
func NewReflectionTransformation() Transformation {
	return NewTransformation([]float64{
		1, 0, 0,
		0, -1, 0,
	})
}

const degreeToRad = math.Pi / float64(180)

// NewRotationTransformation creates a transformation which rotates by the specified amount.
func NewRotationTransformation(degrees float64) Transformation {
	radians := degrees * degreeToRad

	sin, cos := math.Sincos(radians)
	// Because the coordinate space is the opposite to usual maths, the
	// matrix is inverted for the clockwise transform we're doing.
	return NewTransformation([]float64{
		+cos, -sin, 0,
		+sin, +cos, 0,
	})
}

// Apply applies the transformation to a point.
func (t Transformation) Apply(point image.Point) image.Point {
	// See https://en.wikipedia.org/wiki/Matrix_multiplication#Matrix_product_.28two_matrices.29
	// Square matrix and column vector (the point)
	x := float64(point.X)
	y := float64(point.Y)
	z := float64(1)

	x1 := (t.a * x) + (t.b * y) + (t.c * z)
	y1 := (t.p * x) + (t.q * y) + (t.r * z)

	return image.Point{int(round.ToEven(x1, 0)), int(round.ToEven(y1, 0))}
}

// Combine combines two transformations into a single operation.
func (t Transformation) Combine(t2 Transformation) Transformation {
	return NewTransformation([]float64{
		(t.a * t2.a) + (t.b * t2.p) + (t.c * t2.u),
		(t.a * t2.b) + (t.b * t2.q) + (t.c * t2.v),
		(t.a * t2.c) + (t.b * t2.r) + (t.c * t2.w),
		(t.p * t2.a) + (t.q * t2.p) + (t.r * t2.u),
		(t.p * t2.b) + (t.q * t2.q) + (t.r * t2.v),
		(t.p * t2.c) + (t.q * t2.r) + (t.r * t2.w),
		(t.u * t2.a) + (t.v * t2.p) + (t.w * t2.u),
		(t.u * t2.b) + (t.v * t2.q) + (t.w * t2.v),
		(t.u * t2.c) + (t.v * t2.r) + (t.w * t2.w),
	})
}

// Eq compares two transformations against each other.
func (t Transformation) Eq(t2 Transformation) bool {
	return t.a == t2.a &&
		t.b == t2.b &&
		t.c == t2.c &&
		t.p == t2.p &&
		t.q == t2.q &&
		t.r == t2.r &&
		t.u == t2.u &&
		t.v == t2.v &&
		t.w == t2.w
}
