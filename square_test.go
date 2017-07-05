package raster

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/colornames"
)

func TestSquareFunction(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))

	l := NewSquare(image.Point{5, 10}, 20, colornames.White)
	l.Draw(img)

	testColor("top left", t, img.At(5, 10), colornames.White)
	testColor("top right", t, img.At(25, 10), colornames.White)
	testColor("bottom left", t, img.At(5, 30), colornames.White)
	testColor("bottom right", t, img.At(25, 30), colornames.White)

	testColor("outside", t, img.At(2, 2), color.RGBA{})
}

func testColor(name string, t *testing.T, actual, expected color.Color) {
	if actual != expected {
		t.Errorf("%s: expected %v, got %v", name, expected, actual)
	}
}

func TestSquareSize(t *testing.T) {
	s := NewSquare(image.Point{}, 100, colornames.White)
	if s.Bounds().Dx() != 100 {
		t.Errorf("expected a width of 100, but got %v", s.Bounds().Dx())
	}
	if s.Bounds().Dy() != 100 {
		t.Errorf("expected a height of 100, but got %v", s.Bounds().Dy())
	}
}
