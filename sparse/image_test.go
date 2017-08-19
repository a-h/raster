package sparse

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"golang.org/x/image/colornames"
)

func TestThatSparseImplementsTheImageInterface(t *testing.T) {
	var s interface{} = new(Image)
	if _, ok := s.(draw.Image); !ok {
		t.Error("expected sparse/image to implement image/draw")
	}
}

func TestColorModel(t *testing.T) {
	img := NewImage(image.Rect(0, 0, 100, 100))

	if img.ColorModel() != color.RGBAModel {
		t.Errorf("expected RGBAModel, but got %v", img.ColorModel())
	}
}

func TestBounds(t *testing.T) {
	img := NewImage(image.Rect(0, 0, 100, 100))

	if img.Bounds().Dx() != 100 {
		t.Errorf("Expected bounds.dx to be 100, but got %v", img.Bounds().Dx())
	}
	if img.Bounds().Dy() != 100 {
		t.Errorf("Expected bounds.dy to be 100, but got %v", img.Bounds().Dy())
	}
}

func TestDrawnBounds(t *testing.T) {
	img := NewImage(image.Rect(0, 0, 100, 100))

	expected := image.Rect(0, 0, 0, 0)
	if !img.DrawnBounds().Eq(expected) {
		t.Errorf("Before we've written anything expected DrawnBounds of %v, but got %v.", expected, img.DrawnBounds())
	}

	img.Set(49, 50, colornames.White)
	expected = image.Rect(49, 50, 49, 50)
	if !img.DrawnBounds().Eq(expected) {
		t.Errorf("After we've written a single pixel  expected DrawnBounds of %v, but got %v.", expected, img.DrawnBounds())
	}

	img.Set(59, 60, colornames.White)
	expected = image.Rect(49, 50, 59, 60)
	if !img.DrawnBounds().Eq(expected) {
		t.Errorf("After we've written two points  expected DrawnBounds of %v, but got %v.", expected, img.DrawnBounds())
	}
}

func TestSetAndAtFunctions(t *testing.T) {
	img := NewImage(image.Rect(0, 0, 100, 100))

	img.Set(49, 50, colornames.White)
	expected := img.At(49, 50)
	if expected != colornames.White {
		t.Errorf("After setting pixel to white, expected to be able to read it as white, but got %v", img.At(49, 50))
	}

	unwritten := img.At(0, 0)
	empty := color.RGBA{}
	if unwritten != empty {
		t.Errorf("Unset pixels should just return an empty struct, but got %v", unwritten)
	}
}
