package raster

import (
	"image"
	"image/color"
	"testing"
)

func TestLineContainsPoint(t *testing.T) {
	l := NewLine(0, 0, 10, 0)
	if !l.ContainsPoint(image.Point{1, 0}) {
		t.Errorf("Expected point 0,0 to be set")
	}
}

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

		l := NewLine(test.from.X, test.from.Y, test.to.X, test.to.Y)
		l.Draw(img, c)
		set, notSet, setIncorrectly, ok := compare(img, test.expected)
		if !ok {
			t.Errorf("%s: %v was set, %v was not set, %v was set incorrectly", test.name, set, notSet, setIncorrectly)
		}
	}
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

func TestShareSameDirectionFunction(t *testing.T) {
	tests := []struct {
		a        *Line
		b        *Line
		expected bool
	}{
		{
			// Top left to bottom right, then more of a similar direction
			a:        NewLine(0, 0, 5, 5),
			b:        NewLine(5, 5, 6, 6),
			expected: true,
		},
		{
			// Back on itself to the left
			a:        NewLine(0, 0, 5, 5),
			b:        NewLine(5, 5, 1, 5),
			expected: false,
		},
		{
			// Back upwards on itself
			a:        NewLine(0, 0, 5, 5),
			b:        NewLine(5, 5, 5, 1),
			expected: false,
		},
		{
			// Continuing to the right
			a:        NewLine(0, 0, 5, 5),
			b:        NewLine(5, 5, 6, 5),
			expected: true,
		},
		{
			// Line going upwards and diagonal to zero, we're actually carrying on the same vector
			a:        NewLine(5, 5, 1, 1),
			b:        NewLine(1, 1, 0, 0),
			expected: true,
		},
	}

	for _, test := range tests {
		actual := test.a.ShareSameDirection(test.b)

		if actual != test.expected {
			t.Errorf("For a (%v) and b (%v), expected increasing %v, but got %v", test.a, test.b, test.expected, actual)
		}
	}
}
