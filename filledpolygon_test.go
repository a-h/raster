package raster

import "testing"
import "golang.org/x/image/colornames"
import "image"

func TestFilledPolygonSquare(t *testing.T) {
	p := NewFilledPolygon(colornames.White, colornames.White, image.Point{0, 0}, image.Point{1000, 0}, image.Point{1000, 1000}, image.Point{0, 1000})
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	p.Draw(img)

	for y := 0; y < 1000; y++ {
		for x := 0; x < 1000; x++ {
			actualColor := img.At(x, y)
			if actualColor != colornames.White {
				t.Errorf("%v, %v: expected white, but got %v", x, y, actualColor)
			}
		}
	}
}

func TestFilledPolygonBounds(t *testing.T) {
	xOffset, yOffset := 50, 100
	a, b, c, d := image.Point{50 + xOffset, 0 + yOffset}, image.Point{100 + xOffset, 50 + yOffset}, image.Point{50 + xOffset, 100 + yOffset}, image.Point{0 + xOffset, 50 + yOffset}
	p := NewFilledPolygon(colornames.White, colornames.White, a, b, c, d)
	actual := p.Bounds()
	if !actual.Eq(image.Rect(50, 100, 150, 200)) {
		t.Errorf("polygon was not expected size: %v", actual)
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
		p := NewFilledPolygon(lineColor, fillColor, test.points...)
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

func BenchmarkFilledPolygon(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	for i := 0; i < b.N; i++ {
		p := NewFilledPolygon(colornames.White, colornames.Blue, image.Point{0, 0}, image.Point{1000, 0}, image.Point{1000, 1000}, image.Point{0, 1000})
		p.Draw(img)
	}
}
