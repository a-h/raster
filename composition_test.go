package raster

import (
	"image"
	"testing"

	"github.com/a-h/raster/affine"

	"golang.org/x/image/colornames"
)

func TestComposition(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))

	circleInsideSquare := NewComposition(image.Point{250, 250},
		NewSquare(image.Point{0, 0}, 500, colornames.Green),
		NewCircle(image.Point{250, 250}, 250, colornames.Maroon))

	circleInsideSquare.Draw(img)

	// Validate that the square was drawn, by checking each corner.
	// Top left
	if img.At(250, 250) != colornames.Green {
		t.Errorf("{250, 250}: expected green, got %v", img.At(0, 0))
	}
	// Top right
	if img.At(250+500, 250) != colornames.Green {
		t.Errorf("{250+500, 250}: expected green, got %v", img.At(0, 0))
	}
	// Bottom right
	if img.At(250+500, 250+500) != colornames.Green {
		t.Errorf("{250+500, 250+500}: expected green, got %v", img.At(0, 0))
	}
	// Bottom left
	if img.At(250, 250+500) != colornames.Green {
		t.Errorf("{250, 250+500}: expected green, got %v", img.At(0, 0))
	}

	// Validate that the circle was drawn, by checking major points.
	if img.At(250+250, 250) != colornames.Maroon {
		t.Errorf("Top of circle missing.")
	}
	if img.At(250+250, 250+250+250) != colornames.Maroon {
		t.Errorf("Bottom of circle missing.")
	}
	if img.At(250, 250+250) != colornames.Maroon {
		t.Errorf("Left of circle missing.")
	}
	if img.At(250+250+250, 250+250) != colornames.Maroon {
		t.Errorf("Right of circle missing.")
	}
}

func TestCompositionWithAffineTransformation(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 40, 40))

	// Make a diamond and turn it into a square.
	composition := NewComposition(image.Point{20, 10},
		NewPolygon(colornames.White, image.Point{0, 5}, image.Point{5, 0}, image.Point{10, 5}, image.Point{5, 10}))

	composition.Draw(img)

	// We should now have a white square(ish).

	// Validate that the square was drawn, by checking each corner.
	// Top left
	if img.At(20+0, 10+5) != colornames.White {
		t.Errorf("{20, 15}: expected white, got %v", img.At(20, 15))
	}
	if img.At(20+5, 10+0) != colornames.White {
		t.Errorf("{25, 10}: expected white, got %v", img.At(25, 10))
	}
	if img.At(20+10, 10+5) != colornames.White {
		t.Errorf("{30, 15}: expected white, got %v", img.At(30, 15))
	}
	if img.At(20+5, 10+10) != colornames.White {
		t.Errorf("{25, 20}: expected white, got %v", img.At(25, 20))
	}

	// The rotation applies from the orign.

	// Apply the transformation.
	composition.Transformation = affine.NewRotationTransformation(-45)
	img2 := image.NewRGBA(image.Rect(0, 0, 40, 40))
	composition.Draw(img2)

	topLeft := image.Point{24, 6}
	if img2.At(topLeft.X, topLeft.Y) != colornames.White {
		t.Errorf("Top left corner was not in correct position")
	}
	topRight := image.Point{31, 9}
	if img2.At(topRight.X, topRight.Y) != colornames.White {
		t.Errorf("Top right corner was not in correct position")
	}
	bottomLeft := image.Point{24, 14}
	if img2.At(bottomLeft.X, bottomLeft.Y) != colornames.White {
		t.Errorf("Bottom left corner was not in correct position")
	}
	bottomRight := image.Point{31, 14}
	if img2.At(bottomRight.X, bottomRight.Y) != colornames.White {
		t.Errorf("Bottom right corner was not in correct position")
	}
}
