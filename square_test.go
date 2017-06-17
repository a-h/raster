package raster

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func TestSquareFunction(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))

	l := NewSquare(5, 10, 20, colornames.White)
	l.Draw(img)

	testColor("top left", t, img.At(5, 10), colornames.White)
	testColor("top right", t, img.At(25, 10), colornames.White)
	testColor("bottom left", t, img.At(5, 30), colornames.White)
	testColor("bottom right", t, img.At(25, 30), colornames.White)

	testColor("outside", t, img.At(2, 2), color.Transparent)
}

func testColor(name string, t *testing.T, actual, expected color.Color) {
	if actual != expected {
		t.Errorf("%s: expected %v, got %v", name, expected, actual)
	}
}
