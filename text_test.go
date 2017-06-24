package raster

import (
	"image"
	"testing"

	"golang.org/x/image/colornames"
)

func TestText(t *testing.T) {
	tests := []struct {
		toWrite  string
		expected image.Rectangle
	}{
		{
			toWrite: "test",
			// (7 pixels * 4 characters) - the pixel at the end
			expected: image.Rect(0, 0, 7*4-1, 13),
		},
		{
			toWrite: "TEST",
			// (7 pixels * 4 characters) - the pixel at the end
			expected: image.Rect(0, 0, 7*4-1, 13),
		},
	}

	for _, test := range tests {
		//img := image.NewRGBA(image.Rect(0, 0, 300, 300))
		text := NewText(image.Point{0, 0}, test.toWrite, colornames.White)

		bounds := text.Bounds()
		if bounds.Dx() != test.expected.Dx() {
			t.Errorf("expected '%s' to be %v wide, but was %v", test.toWrite, test.expected.Dx(), bounds.Dx())
		}
		if bounds.Dy() != test.expected.Dy() {
			t.Errorf("expected '%s' to be %v tall, but was %v", test.toWrite, test.expected.Dy(), bounds.Dy())
		}
	}
}
