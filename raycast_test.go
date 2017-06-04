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

func TestCalculateDirection(t *testing.T) {
	tests := []struct {
		name     string
		l1       *Line
		l2       *Line
		expected Direction
	}{
		{
			name:     "horizontal",
			l1:       NewLine(0, 0, 10, 0),
			l2:       NewLine(10, 0, 10, 0),
			expected: None,
		},
		{
			name:     "right arrow",
			l1:       NewLine(0, 0, 5, 5),
			l2:       NewLine(5, 5, 0, 10),
			expected: Right,
		},
		{
			name:     "left arrow",
			l1:       NewLine(5, 0, 0, 5),
			l2:       NewLine(0, 5, 10, 10),
			expected: Left,
		},
		{
			name:     "up arrow",
			l1:       NewLine(0, 10, 5, 0),
			l2:       NewLine(5, 0, 10, 10),
			expected: Up,
		},
		{
			name:     "vertical",
			l1:       NewLine(0, 0, 0, 10),
			l2:       NewLine(0, 10, 0, 20),
			expected: None,
		},
	}

	for _, test := range tests {
		actual := CalculateDirection(test.l1, test.l2)
		if actual != test.expected {
			t.Errorf("%s: expected %v, but got %v", test.name, test.expected, actual)
		}
	}
}
