package affine

import "testing"

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
