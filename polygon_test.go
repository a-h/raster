package raster

import (
	"image"
	"testing"
)

func TestBounds(t *testing.T) {
	tests := []struct {
		p        Polygon
		expected image.Rectangle
	}{
		{
			// 10 x 10 square
			p:        NewPolygon(image.Point{}, image.Point{10, 0}, image.Point{10, 10}, image.Point{0, 10}),
			expected: image.Rect(0, 0, 10, 10),
		},
		{
			// 5 x 5 square
			p:        NewPolygon(image.Point{5, 5}, image.Point{10, 5}, image.Point{10, 10}, image.Point{5, 10}),
			expected: image.Rect(0, 0, 5, 5),
		},
	}

	for _, test := range tests {
		actual := test.p.Bounds()
		if !actual.Eq(test.expected) {
			t.Errorf("For polygon %v, expected bounds %v, but got %v", test.p, test.expected, actual)
		}
	}
}
