package raster

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func TestText(t *testing.T) {
	tests := []struct {
		toWrite  string
		expected [][]int
	}{
		{
			toWrite: "H",
			expected: [][]int{
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
				[]int{1, 1, 1, 1, 1, 1, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
			},
		},
		{
			toWrite: "e",
			expected: [][]int{
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0, 0},
				[]int{0, 1, 1, 1, 1, 0, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
				[]int{1, 1, 1, 1, 1, 1, 0},
				[]int{1, 0, 0, 0, 0, 0, 0},
				[]int{1, 0, 0, 0, 0, 1, 0},
				[]int{0, 1, 1, 1, 1, 0, 0},
			},
		},
	}

	for _, test := range tests {
		text := NewText(image.Point{0, 0}, test.toWrite, colornames.White)
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		text.Draw(img)

		for y, yi := range test.expected {
			for x, xi := range yi {
				actual := img.At(x, y)
				shouldBeWhite := (xi == 1)
				expected := color.RGBA{}
				if shouldBeWhite {
					expected = colornames.White
				}
				if actual != expected {
					t.Errorf("for '%v' at {%v, %v}: expected %v, but got %v", test.toWrite, x, y, expected, actual)
				}
			}
		}
	}
}

func TestTextBounds(t *testing.T) {
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
