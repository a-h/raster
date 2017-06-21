package affine

import (
	"image"
	"testing"
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
	}

	for _, test := range tests {
		transformation := NewRotationTransformation(test.degrees)
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