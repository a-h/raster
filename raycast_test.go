package raster

import (
	"image"
	"testing"
)

func TestRaycastPolygon(t *testing.T) {
	// Form a diamond.
	a := image.Point{5, 0}
	b := image.Point{10, 5}
	c := image.Point{5, 10}
	d := image.Point{0, 5}

	p := NewPolygon(a, b, c, d)

	tests := []struct {
		point              image.Point
		expectedBoundaries int
	}{
		{
			point:              image.Point{0, 0},
			expectedBoundaries: 0,
		},
		{
			point:              image.Point{4, 4},
			expectedBoundaries: 1,
		},
		{
			point:              image.Point{5, 5},
			expectedBoundaries: 1,
		},
		{
			point:              image.Point{7, 10},
			expectedBoundaries: 2,
		},
	}

	for _, test := range tests {
		actualBoundaries := ScanLine(test.point.Y, p)

		if actualBoundaries[test.point.X] != test.expectedBoundaries {
			t.Errorf("for %v, expected %d boundaries, but got %d", test.point, test.expectedBoundaries, actualBoundaries)
		}
	}
}
