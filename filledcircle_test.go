package raster

import (
	"image"
	"testing"

	"golang.org/x/image/colornames"
)

func TestFilledCircle(t *testing.T) {
	c := NewFilledCircle(image.Point{50, 50}, 50, colornames.White, colornames.Aliceblue)

	img := image.NewRGBA(image.Rect(0, 0, 101, 101))
	c.Draw(img)

	// Check the basic shape is there.
	if img.At(50, 50) == colornames.White {
		t.Error("expected nothing in the center")
	}

	if img.At(0, 0) == colornames.White {
		t.Error("expected nothing at the top left")
	}

	if img.At(0, 50) != colornames.White {
		t.Error("expected the left edge of the circle to be set")
	}

	if img.At(100, 50) != colornames.White {
		t.Error("expected the right edge of the circle to be set")
	}

	if img.At(50, 0) != colornames.White {
		t.Error("expected the top edge of the circle to be set")
	}

	if img.At(50, 100) != colornames.White {
		t.Error("expected the bottom edge of the circle to be set")
	}

	if img.At(50, 50) != colornames.Aliceblue {
		t.Error("expected the middke the circle to be filled")
	}
}

func TestFilledCircleBounds(t *testing.T) {
	radius := 1000
	c := NewFilledCircle(image.Point{100, 100}, radius, colornames.White, colornames.Aliceblue)

	if c.Bounds().Dx() != radius*2 {
		t.Errorf("expected 2000 diameter, but got %v", c.Bounds().Dx())
	}

	if c.Bounds().Dy() != radius*2 {
		t.Errorf("expected 2000 diameter, but got %v", c.Bounds().Dy())
	}
}

func BenchmarkFilledCircle(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	for i := 0; i < b.N; i++ {
		p := NewFilledCircle(image.Point{500, 500}, 500, colornames.White, colornames.Aliceblue)
		p.Draw(img)
	}
}
