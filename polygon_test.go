package raster

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func TestBounds(t *testing.T) {
	tests := []struct {
		p        Polygon
		expected image.Rectangle
	}{
		{
			// 10 x 10 square
			p:        NewPolygon(colornames.White, image.Point{}, image.Point{10, 0}, image.Point{10, 10}, image.Point{0, 10}),
			expected: image.Rect(0, 0, 10, 10),
		},
		{
			// 5 x 5 square
			p:        NewPolygon(colornames.White, image.Point{5, 5}, image.Point{10, 5}, image.Point{10, 10}, image.Point{5, 10}),
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

func TestPolygon(t *testing.T) {
	tests := []struct {
		name          string
		size          image.Rectangle
		points        []image.Point
		expectedDrawn []image.Point
	}{
		{
			name:          "line",
			size:          image.Rect(0, 0, 3, 3),
			points:        []image.Point{image.Point{0, 0}, image.Point{0, 2}},
			expectedDrawn: []image.Point{image.Point{0, 0}, image.Point{0, 1}, image.Point{0, 2}},
		},
		{
			name:          "square around the edges",
			size:          image.Rect(0, 0, 3, 3),
			points:        []image.Point{image.Point{0, 0}, image.Point{2, 0}, image.Point{2, 2}, image.Point{0, 2}},
			expectedDrawn: []image.Point{image.Point{0, 0}, image.Point{1, 0}, image.Point{2, 0}, image.Point{0, 1}, image.Point{2, 1}, image.Point{0, 2}, image.Point{1, 2}, image.Point{2, 2}},
		},
	}

	lineColor := colornames.Red

	for _, test := range tests {
		img := image.NewRGBA(test.size)
		p := NewPolygon(lineColor, test.points...)
		p.Draw(img)

		colors := mapColors(img)

		actualDrawn := colors[lineColor]
		drawnAccidentally, notDrawn, _ := discover(actualDrawn, test.expectedDrawn)
		if len(drawnAccidentally) > 0 {
			t.Errorf("%s: %v should not have been drawn", test.name, drawnAccidentally)
		}
		if len(notDrawn) > 0 {
			t.Errorf("%s: %v was missing points %v", test.name, actualDrawn, notDrawn)
		}
	}
}

func discover(a []image.Point, b []image.Point) (onlyInA []image.Point, onlyInB []image.Point, inBoth []image.Point) {
	both := make(map[image.Point]interface{})

	for _, ap := range a {
		found := false
		for _, bp := range b {
			if ap == bp {
				both[ap] = true
				found = true
				break
			}
		}
		if !found {
			onlyInA = append(onlyInA, ap)
		}
	}

	inBoth = make([]image.Point, len(both))
	i := 0
	for k := range both {
		inBoth[i] = k
		i++
	}

	return onlyInA, onlyInB, inBoth
}

func mapColors(img image.Image) map[color.Color][]image.Point {
	rv := make(map[color.Color][]image.Point)

	for x := 0; x <= img.Bounds().Dx(); x++ {
		for y := 0; y <= img.Bounds().Dy(); y++ {
			c := img.At(x, y)
			points, ok := rv[c]
			if !ok {
				points = []image.Point{}
			}
			rv[c] = append(points, image.Point{x, y})
		}
	}

	return rv
}
