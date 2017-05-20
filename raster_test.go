package raster

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func TestDrawLines(t *testing.T) {
	tests := []struct {
		name     string
		size     image.Rectangle
		from, to image.Point
		expected []image.Point
	}{
		{
			name:     "A point",
			size:     image.Rect(0, 0, 1, 1),
			from:     image.Point{1, 1},
			to:       image.Point{1, 1},
			expected: []image.Point{image.Point{1, 1}},
		},
		{
			name:     "2px horizontal line (L-R)",
			size:     image.Rect(0, 0, 10, 1),
			from:     image.Point{0, 0},
			to:       image.Point{2, 0},
			expected: []image.Point{image.Point{0, 0}, image.Point{1, 0}, image.Point{2, 0}},
		},
		{
			name:     "2px horizontal line (R-L)",
			size:     image.Rect(0, 0, 10, 1),
			from:     image.Point{2, 0},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{0, 0}, image.Point{1, 0}, image.Point{2, 0}},
		},
		{
			name:     "2px vertical line (top-bottom)",
			size:     image.Rect(0, 0, 10, 0),
			from:     image.Point{0, 0},
			to:       image.Point{0, 2},
			expected: []image.Point{image.Point{0, 0}, image.Point{0, 1}, image.Point{0, 2}},
		},
		{
			name:     "2px vertical line (bottom-top)",
			size:     image.Rect(0, 0, 10, 0),
			from:     image.Point{0, 2},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{0, 0}, image.Point{0, 1}, image.Point{0, 2}},
		},
		{
			name:     "45 degrees",
			size:     image.Rect(0, 0, 2, 2),
			from:     image.Point{0, 2},
			to:       image.Point{2, 0},
			expected: []image.Point{image.Point{0, 2}, image.Point{1, 1}, image.Point{2, 0}},
		},
		{
			name:     "45 degrees (2)",
			size:     image.Rect(0, 0, 4, 4),
			from:     image.Point{2, 2},
			to:       image.Point{4, 0},
			expected: []image.Point{image.Point{2, 2}, image.Point{3, 1}, image.Point{4, 0}},
		},
		{
			name:     "135 degrees",
			size:     image.Rect(0, 0, 2, 2),
			from:     image.Point{0, 0},
			to:       image.Point{2, 2},
			expected: []image.Point{image.Point{0, 0}, image.Point{1, 1}, image.Point{2, 2}},
		},
		{
			name:     "225 degrees",
			size:     image.Rect(0, 0, 2, 2),
			from:     image.Point{2, 0},
			to:       image.Point{0, 2},
			expected: []image.Point{image.Point{0, 2}, image.Point{1, 1}, image.Point{2, 2}},
		},
		{
			name:     "270 degrees",
			size:     image.Rect(0, 0, 2, 2),
			from:     image.Point{2, 2},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{2, 2}, image.Point{1, 1}, image.Point{0, 0}},
		},
		{
			name:     "270 degrees (part)",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{3, 3},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{3, 3}, image.Point{2, 2}, image.Point{1, 1}, image.Point{0, 0}},
		},
		{
			name:     "bottom to top",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{0, 2},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{0, 2}, image.Point{0, 1}, image.Point{0, 0}},
		},
	}

	c := color.RGBA{R: 0xFF, G: 0x55, B: 0x55}

	for _, test := range tests {
		img := image.NewRGBA(test.size)
		DrawLine(img, test.from.X, test.from.Y, test.to.X, test.to.Y, c)
		set, notSet, setIncorrectly, ok := compare(img, test.expected)
		if !ok {
			t.Errorf("%s: %v was set, %v was not set, %v was set incorrectly", test.name, set, notSet, setIncorrectly)
		}
	}
}

func TestFilledPolygon(t *testing.T) {
	tests := []struct {
		name           string
		size           image.Rectangle
		points         []image.Point
		expectedDrawn  []image.Point
		expectedFilled []image.Point
	}{
		{
			name:           "simple square",
			size:           image.Rect(0, 0, 4, 4),
			points:         []image.Point{image.Point{0, 0}, image.Point{2, 0}, image.Point{2, 2}, image.Point{0, 2}},
			expectedDrawn:  []image.Point{image.Point{0, 0}, image.Point{1, 0}, image.Point{2, 0}, image.Point{2, 1}, image.Point{2, 2}, image.Point{1, 2}, image.Point{0, 2}, image.Point{0, 1}},
			expectedFilled: []image.Point{image.Point{1, 1}},
		},
		{
			name:           "square surrounded by space",
			size:           image.Rect(0, 0, 5, 5),
			points:         []image.Point{image.Point{1, 1}, image.Point{3, 1}, image.Point{3, 3}, image.Point{1, 3}},
			expectedDrawn:  []image.Point{image.Point{1, 1}, image.Point{2, 1}, image.Point{3, 1}, image.Point{3, 2}, image.Point{3, 3}, image.Point{2, 3}, image.Point{1, 3}, image.Point{1, 2}},
			expectedFilled: []image.Point{image.Point{2, 2}},
		},
		{
			name:           "simple unfilled triangle",
			size:           image.Rect(0, 0, 5, 5),
			points:         []image.Point{image.Point{0, 0}, image.Point{2, 2}, image.Point{0, 2}},
			expectedDrawn:  []image.Point{image.Point{0, 0}, image.Point{1, 1}, image.Point{2, 2}, image.Point{1, 2}, image.Point{0, 2}, image.Point{0, 1}},
			expectedFilled: []image.Point{},
		},
		{
			name:           "simple filled triangle",
			size:           image.Rect(0, 0, 5, 5),
			points:         []image.Point{image.Point{0, 0}, image.Point{3, 3}, image.Point{0, 3}},
			expectedDrawn:  []image.Point{image.Point{0, 0}, image.Point{1, 1}, image.Point{2, 2}, image.Point{3, 3}, image.Point{2, 3}, image.Point{1, 3}, image.Point{0, 3}, image.Point{0, 2}, image.Point{0, 1}},
			expectedFilled: []image.Point{image.Point{1, 2}},
		},
	}

	fillColor := colornames.White
	lineColor := colornames.Red

	for _, test := range tests {
		img := image.NewRGBA(test.size)
		DrawFilledPolygon(img, lineColor, fillColor, test.points...)

		colors := mapColors(img)

		actualDrawn := colors[lineColor]
		drawnAccidentally, notDrawn, _ := discover(actualDrawn, test.expectedDrawn)
		if len(drawnAccidentally) > 0 {
			t.Errorf("%s: %v should not have been drawn", test.name, drawnAccidentally)
		}
		if len(notDrawn) > 0 {
			t.Errorf("%s: %v was missing points %v", test.name, actualDrawn, notDrawn)
		}

		actualFilled := colors[fillColor]
		filledAccidentally, notFilled, _ := discover(actualFilled, test.expectedFilled)
		if len(drawnAccidentally) > 0 {
			t.Errorf("%s: %v should not have been filled", test.name, filledAccidentally)
		}
		if len(notFilled) > 0 {
			t.Errorf("%s: %v was missing filled points %v", test.name, actualFilled, notFilled)
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
		DrawPolygon(img, lineColor, test.points...)

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

func compare(img *image.RGBA, activePixels []image.Point) (set, notSet, setIncorrectly []image.Point, ok bool) {
	// Make a map of points to speed up comparison instead of sorting them.
	expectedPoints := make(map[image.Point]interface{}, len(activePixels))
	for _, p := range activePixels {
		expectedPoints[p] = true
	}

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			p := image.Point{x, y}
			_, expected := expectedPoints[p]
			isSet := img.At(x, y) != color.RGBA{}

			if isSet {
				set = append(set, p)
			}

			if expected {
				if !isSet {
					notSet = append(notSet, p)
				}
			} else {
				if isSet {
					setIncorrectly = append(setIncorrectly, p)
				}
			}
		}
	}

	return set, notSet, setIncorrectly, ((len(notSet) == 0) && (len(setIncorrectly) == 0))
}
