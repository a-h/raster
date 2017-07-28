package affine

import (
	"fmt"
	"image"
	"math"
	"reflect"
	"testing"

	"github.com/a-h/linear/tolerance"
)

func TestMatrixCombination(t *testing.T) {
	a := NewTransformation([]float64{
		-1, 4, -6,
		8, 5, 16,
		2, 8, 5,
	})

	b := NewTransformation([]float64{
		12, 7, 6,
		8, 0, 5,
		3, 2, 4,
	})

	expected := NewTransformation([]float64{
		-8, -19, -10,
		184, 88, 137,
		103, 24, 72,
	})

	actual := a.Combine(b)

	if actual.Eq(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestRotationTransformation(t *testing.T) {
	// A 90 degree transformation around the origin should result
	// in:
	// 0,0 being translated to 0,0
	// 10,0 being translated to 0,10
	tests := []struct {
		input    image.Point
		expected image.Point
		degrees  float64
	}{
		{
			input:    image.Point{0, 0},
			expected: image.Point{0, 0},
			degrees:  90,
		},
		{
			input:    image.Point{10, 0},
			expected: image.Point{0, 10},
			degrees:  90,
		},
		{
			input:    image.Point{10, 0},
			expected: image.Point{-10, 0},
			degrees:  180,
		},
		{
			input:    image.Point{10, 0},
			expected: image.Point{0, -10},
			degrees:  270,
		},
		{
			input:    image.Point{5, 5},
			expected: image.Point{0, 7},
			degrees:  45,
		},
		{
			input:    image.Point{5, 5},
			expected: image.Point{7, 0},
			degrees:  -45,
		},
	}

	for _, test := range tests {
		transformation := NewRotationTransformation(test.degrees)
		actual := transformation.Apply(test.input)
		if !actual.Eq(test.expected) {
			t.Errorf("for input %v rotated by %v, expected %v, got %v", test.input, test.degrees, test.expected, actual)
		}
	}
}

func convertDegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / float64(180))
}

func TestRotationAlgebra(t *testing.T) {
	// With a triangle with base and height of 5.
	// We have a right angle, 90 degrees. Then the two remaining angles are: (180-90)/2 = 45 degrees
	// Given that the interior angle is 45 degrees, then it's 45 degrees to the X axis, or Y axis.

	// The length of the hypotenuse of this triangle is:
	// sqrt((5*5)+(5*5))=7.071067811865475

	// If we rotate anticlockwise around the origin, the result should be:
	// 0, 7.071067811865475
	// Clockwise, should be:
	// 7.071067811865475, 0

	// To rotate the point, it's cos(t), sin(t)
	// In algebra, clockwise rotation is given by negative numbers.
	theta := convertDegreesToRadians(float64(45))
	fmt.Printf("Theta in rads: %v\n", theta)

	// x' = x cos f - y sin f
	// y' = y cos f + x sin f
	x := float64(5)
	y := float64(5)

	sin45, cos45 := math.Sincos(theta)
	expectedSin45 := 0.7071067811865476
	if sin45 != expectedSin45 {
		t.Errorf("sin(45) expected %v, but got %v\n", expectedSin45, sin45)
	}

	expectedCos45 := 0.7071067811865475
	if cos45 != expectedCos45 {
		t.Errorf("cos(45) expected %v, but got %v\n", expectedCos45, cos45)
	}

	xp := (x * math.Cos(theta)) - (y * math.Sin(theta))
	yp := (y * math.Cos(theta)) + (y * math.Sin(theta))

	expectedX := float64(0)
	expectedY := 7.071067811865475

	fmt.Println(xp, yp)
	if !tolerance.IsWithin(xp, expectedX, tolerance.ThreeDecimalPlaces) {
		t.Errorf("x: expected %v, but got %v", expectedX, xp)
	}
	if !tolerance.IsWithin(yp, expectedY, tolerance.ThreeDecimalPlaces) {
		t.Errorf("y: expected %v, but got %v", expectedY, yp)
	}
}

func TestTranslationTransformation(t *testing.T) {
	// A 90 degree transformation around the origin should result
	// in:
	// 0,0 being translated to 0,0
	// 10,0 being translated to 0,10
	tests := []struct {
		input    image.Point
		expected image.Point
		x        int
		y        int
	}{
		{
			input:    image.Point{0, 0},
			x:        0,
			y:        0,
			expected: image.Point{0, 0},
		},
		{
			input:    image.Point{0, 0},
			x:        10,
			y:        0,
			expected: image.Point{10, 0},
		},
		{
			input:    image.Point{0, 0},
			x:        0,
			y:        10,
			expected: image.Point{0, 10},
		},
		{
			input:    image.Point{0, 0},
			x:        10,
			y:        10,
			expected: image.Point{10, 10},
		},
		{
			input:    image.Point{0, 0},
			x:        -10,
			y:        -10,
			expected: image.Point{-10, -10},
		},
	}

	for _, test := range tests {
		transformation := NewTranslationTransformation(test.x, test.y)
		actual := transformation.Apply(test.input)
		if !actual.Eq(test.expected) {
			t.Errorf("expected %v, got %v", test.expected, actual)
		}
	}
}
func TestTransformationEquality(t *testing.T) {
	t1 := NewTransformation([]float64{0, 0, 0, 0, 0, 0, 0, 0, 0})
	t2 := NewTransformation([]float64{0, 0, 0, 0, 0, 0, 0, 0, 0})

	if !t1.Eq(t2) {
		t.Errorf("Expected Eq to return true because the two transformation matrices are identical")
	}

	t2 = NewTransformation([]float64{0, 0, 1, 0, 0, 0})

	if t1.Eq(t2) {
		t.Errorf("Expected Eq to return false")
	}
}

func TestScaleTransformation(t *testing.T) {
	tests := []struct {
		name        string
		input       []image.Point
		scaleWidth  float64
		scaleHeight float64
		expected    []image.Point
	}{
		{
			name: "Do nothing",
			input: []image.Point{
				image.Point{0, 0},
				image.Point{1, 0},
			},
			scaleWidth:  1,
			scaleHeight: 1,
			expected: []image.Point{
				image.Point{0, 0},
				image.Point{1, 0},
			},
		},
		{
			name: "Half the width",
			input: []image.Point{
				image.Point{0, 0},
				image.Point{1, 0},
				image.Point{2, 0},
				image.Point{3, 0},
				image.Point{4, 0},
			},
			scaleWidth:  0.5,
			scaleHeight: 1,
			expected: []image.Point{
				image.Point{0, 0},
				image.Point{0, 0},
				image.Point{1, 0},
				image.Point{2, 0},
				image.Point{2, 0},
			},
		},
		{
			name: "Ignore high values",
			input: []image.Point{
				image.Point{0, 0},
				image.Point{1, 0},
				image.Point{2, 0},
				image.Point{3, 0},
				image.Point{4, 0},
			},
			scaleWidth:  1.5,
			scaleHeight: 1.5,
			expected: []image.Point{
				image.Point{0, 0},
				image.Point{1, 0},
				image.Point{2, 0},
				image.Point{3, 0},
				image.Point{4, 0},
			},
		},
	}

	for _, test := range tests {
		transformation := NewScaleTransformation(test.scaleWidth, test.scaleHeight)
		actual := []image.Point{}

		for _, p := range test.input {
			actual = append(actual, transformation.Apply(p))
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("%s: expected %v, got %v", test.name, test.expected, actual)
		}
	}
}

func TestIdentityMatrixTransform(t *testing.T) {
	tests := []struct {
		input []image.Point
	}{
		{
			input: []image.Point{
				image.Point{0, 0},
				image.Point{1, 0},
			},
		},
		{
			input: []image.Point{
				image.Point{0, 0},
				image.Point{1, 0},
				image.Point{2, 0},
				image.Point{3, 0},
				image.Point{4, 0},
			},
		},
		{
			input: []image.Point{
				image.Point{0, 0},
				image.Point{1, 1},
				image.Point{2, 2},
				image.Point{3, 3},
				image.Point{4, 4},
			},
		},
		{
			input: []image.Point{
				image.Point{0, 0},
				image.Point{0, 1},
				image.Point{0, 2},
				image.Point{0, 3},
				image.Point{0, 4},
			},
		},
	}

	for _, test := range tests {
		transformation := NewTransformation(IdentityMatrix)
		actual := []image.Point{}

		for _, p := range test.input {
			actual = append(actual, transformation.Apply(p))
		}

		if !reflect.DeepEqual(actual, test.input) {
			t.Errorf("for input %v, got %v", test.input, actual)
		}
	}
}
