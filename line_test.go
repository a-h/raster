package raster

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func TestLineContainsPointPositive(t *testing.T) {
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
			size:     image.Rect(0, 0, 3, 3),
			from:     image.Point{2, 0},
			to:       image.Point{0, 2},
			expected: []image.Point{image.Point{2, 0}, image.Point{1, 1}, image.Point{0, 2}},
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
		{
			name:     "right to left upwards diagonal",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{3, 3},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{3, 3}, image.Point{2, 2}, image.Point{1, 1}, image.Point{0, 0}},
		},
		{
			name:     "left to right upwards diagonal",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{0, 3},
			to:       image.Point{3, 0},
			expected: []image.Point{image.Point{0, 3}, image.Point{1, 2}, image.Point{2, 1}, image.Point{3, 0}},
		},
		{
			name:     "left to right downwards diagonal",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{3, 0},
			to:       image.Point{0, 3},
			expected: []image.Point{image.Point{0, 3}, image.Point{1, 2}, image.Point{2, 1}, image.Point{3, 0}},
		},
		{
			name:     "right to left downwards diagonal",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{0, 0},
			to:       image.Point{3, 3},
			expected: []image.Point{image.Point{0, 0}, image.Point{1, 1}, image.Point{2, 2}, image.Point{3, 3}},
		},
		{
			name:     "sharp slope in y",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{0, 0},
			to:       image.Point{1, 6},
			expected: []image.Point{image.Point{0, 0}, image.Point{0, 1}, image.Point{0, 2}, image.Point{0, 3}, image.Point{0, 4}, image.Point{0, 5}, image.Point{1, 6}},
		},
		{
			name:     "sharp slope in y (reverse)",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{1, 6},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{0, 0}, image.Point{0, 1}, image.Point{0, 2}, image.Point{0, 3}, image.Point{0, 4}, image.Point{0, 5}, image.Point{1, 6}},
		},
		{
			name:     "sharp slope in x",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{0, 0},
			to:       image.Point{6, 1},
			expected: []image.Point{image.Point{0, 0}, image.Point{1, 0}, image.Point{2, 0}, image.Point{3, 0}, image.Point{4, 0}, image.Point{5, 0}, image.Point{6, 1}},
		},
		{
			name:     "sharp slope in x (reversed)",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{6, 1},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{0, 0}, image.Point{1, 0}, image.Point{2, 0}, image.Point{3, 0}, image.Point{4, 0}, image.Point{5, 0}, image.Point{6, 1}},
		},
	}

	for _, test := range tests {
		l := NewLine(test.from, test.to, colornames.White)
		for _, e := range test.expected {
			if !l.ContainsPoint(e) {
				t.Errorf("%s: expected point %v to be set, but it wasn't", test.name, e)
			}
		}
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
		{
			name:     "right to left upwards diagonal",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{3, 3},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{3, 3}, image.Point{2, 2}, image.Point{1, 1}, image.Point{0, 0}},
		},
		{
			name:     "left to right upwards diagonal",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{0, 3},
			to:       image.Point{3, 0},
			expected: []image.Point{image.Point{0, 3}, image.Point{1, 2}, image.Point{2, 1}, image.Point{3, 0}},
		},
		{
			name:     "left to right downwards diagonal",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{3, 0},
			to:       image.Point{0, 3},
			expected: []image.Point{image.Point{0, 3}, image.Point{1, 2}, image.Point{2, 1}, image.Point{3, 0}},
		},
		{
			name:     "right to left downwards diagonal",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{0, 0},
			to:       image.Point{3, 3},
			expected: []image.Point{image.Point{0, 0}, image.Point{1, 1}, image.Point{2, 2}, image.Point{3, 3}},
		},
		{
			name:     "sharp slope in y",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{0, 0},
			to:       image.Point{1, 6},
			expected: []image.Point{image.Point{0, 0}, image.Point{0, 1}, image.Point{0, 2}, image.Point{0, 3}, image.Point{0, 4}, image.Point{0, 5}, image.Point{1, 6}},
		},
		{
			name:     "sharp slope in y (reverse)",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{1, 6},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{0, 0}, image.Point{0, 1}, image.Point{0, 2}, image.Point{0, 3}, image.Point{0, 4}, image.Point{0, 5}, image.Point{1, 6}},
		},
		{
			name:     "sharp slope in x",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{0, 0},
			to:       image.Point{6, 1},
			expected: []image.Point{image.Point{0, 0}, image.Point{1, 0}, image.Point{2, 0}, image.Point{3, 0}, image.Point{4, 0}, image.Point{5, 0}, image.Point{6, 1}},
		},
		{
			name:     "sharp slope in x (reversed)",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{6, 1},
			to:       image.Point{0, 0},
			expected: []image.Point{image.Point{0, 0}, image.Point{1, 0}, image.Point{2, 0}, image.Point{3, 0}, image.Point{4, 0}, image.Point{5, 0}, image.Point{6, 1}},
		},
		{
			name:     "sharp slope in y (up-left)",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{6, 6},
			to:       image.Point{5, 0},
			expected: []image.Point{image.Point{5, 0}, image.Point{5, 1}, image.Point{5, 2}, image.Point{5, 3}, image.Point{5, 4}, image.Point{5, 5}, image.Point{6, 6}},
		},
		{
			name:     "sharp slope in y (extra)",
			size:     image.Rect(0, 0, 10, 10),
			from:     image.Point{6, 6},
			to:       image.Point{7, 0},
			expected: []image.Point{image.Point{6, 1}, image.Point{6, 2}, image.Point{6, 3}, image.Point{6, 4}, image.Point{6, 5}, image.Point{7, 0}, image.Point{6, 6}},
		},
	}

	for _, test := range tests {
		img := image.NewRGBA(test.size)

		l := NewLine(test.from, test.to, colornames.White)
		l.Draw(img)
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

func TestLineStringFunction(t *testing.T) {
	l := NewLine(image.Point{0, 1}, image.Point{2, 3}, colornames.White)
	if l.String() != "{0, 1} to {2, 3}" {
		t.Errorf("unexpected line string representation")
	}
}

func TestLineBoundsFunction(t *testing.T) {
	tests := []struct {
		from     image.Point
		to       image.Point
		expected image.Rectangle
	}{
		{
			from:     image.Point{0, 0},
			to:       image.Point{100, 100},
			expected: image.Rect(0, 0, 100, 100),
		},
		{
			from:     image.Point{100, 100},
			to:       image.Point{100, 100},
			expected: image.Rect(0, 0, 0, 0),
		},
		{
			from:     image.Point{100, 100},
			to:       image.Point{500, 500},
			expected: image.Rect(0, 0, 400, 400),
		},
		{
			from:     image.Point{10, 10},
			to:       image.Point{0, 0},
			expected: image.Rect(0, 0, 10, 10),
		},
	}

	for _, test := range tests {
		l := NewLine(test.from, test.to, colornames.White)
		actual := l.Bounds()
		if !actual.Eq(test.expected) {
			t.Errorf("For line %v, expected bounds %v, but got %v", l.String(), test.expected, actual)
		}
	}
}

func BenchmarkLine(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	for i := 0; i < b.N; i++ {
		l := NewLine(image.Point{0, 0}, image.Point{1000, 1000}, colornames.White)
		l.Draw(img)
	}
}

func BenchmarkLineContainsPoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := NewLine(image.Point{0, 0}, image.Point{1000, 1000}, colornames.White)
		l.ContainsPoint(image.Point{0, 0})
	}
}
