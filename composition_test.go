package raster

import (
	"image"
	"testing"

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
