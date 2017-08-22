package stage

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/a-h/raster/sparse"

	"golang.org/x/image/colornames"
)

func TestThatStageImplementsTheImageInterface(t *testing.T) {
	var s interface{} = new(Stage)
	if _, ok := s.(draw.Image); !ok {
		t.Error("Expected stage to implement image/draw")
	}
}

func TestBounds(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 250, 250))
	stg := New(img)
	expected := image.Rect(0, 0, 250, 250)
	if stg.Bounds() != expected {
		t.Errorf("Expected bounds of %v, but got %v", expected, stg.Bounds())
	}
}

func TestColorModel(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 250, 250))
	stg := New(img)
	expected := img.ColorModel()
	if stg.ColorModel() != expected {
		t.Errorf("Expected ColorModel of %v, but got %v", expected, stg.ColorModel())
	}
}

func TestThatTheCurrentAndNextFramesAreTheSameSizeAsTheBackdrop(t *testing.T) {
	bounds := image.Rect(0, 0, 250, 250)
	img := image.NewRGBA(bounds)
	stg := New(img)

	stg.ResetCurrentFrame()

	if stg.CurrentFrame.Bounds() != bounds {
		t.Errorf("Expected current frame bounds of %v, but got %v", bounds, stg.CurrentFrame.Bounds())
	}

	stg.ResetNextFrame()

	if stg.NextFrame.Bounds() != bounds {
		t.Errorf("Expected next frame bounds of %v, but got %v", bounds, stg.NextFrame.Bounds())
	}
}

func TestThatItsPossibleToDrawOnABackground(t *testing.T) {
	bounds := image.Rect(0, 0, 250, 250)

	// Create a white background.
	stgBackdrop := image.NewRGBA(bounds)
	draw.Draw(stgBackdrop, stgBackdrop.Bounds(), &image.Uniform{colornames.White}, image.ZP, draw.Src)
	stg := New(stgBackdrop)

	target := image.NewRGBA(bounds)

	// Draw once to set the background.
	stg.Draw(target)
	if stg.At(50, 50) != colornames.White {
		t.Errorf("At (50, 50), before it's set, the pixel should be the background color White, but got %v", stg.At(50, 50))
	}
	stg.NextFrame.Set(50, 50, colornames.Red)
	stg.Draw(target)
	if stg.At(50, 50) != colornames.Red {
		t.Errorf("At (50, 50), after it's set, the pixel should be Red, but got %v", stg.At(50, 50))
	}
	stg.NextFrame.Set(50, 50, colornames.Orange)
	stg.Draw(target)
	if stg.At(50, 50) != colornames.Orange {
		t.Errorf("At (50, 50), after it's set again, the pixel should be Orange, but got %v", stg.At(50, 50))
	}
}

func TestModifyingTheFrameManually(t *testing.T) {
	bounds := image.Rect(0, 0, 250, 250)

	stg := New(image.NewRGBA(bounds))
	stg.Set(50, 50, colornames.Red)

	if stg.At(50, 50) != colornames.Red {
		t.Errorf("It should be possible to modify the current frame, even if it makes no sense to do so")
	}

	stg.Draw(image.NewRGBA(bounds))

	empty := color.RGBA{}
	if stg.At(50, 50) != empty {
		t.Errorf("After drawing, the current frame should be replaced with the new frame")
	}
}

func TestAnimationDeltas(t *testing.T) {
	r := colornames.Red
	g := colornames.Green
	// b := colornames.Blue
	x := color.RGBA{}

	tests := []struct {
		name        string
		firstFrame  []color.Color
		secondFrame []color.Color
	}{
		{
			name: "move red",
			firstFrame: []color.Color{
				r, r, r, x, x,
			},
			secondFrame: []color.Color{
				x, x, r, r, r,
			},
		},
		{
			name: "switch red to green",
			firstFrame: []color.Color{
				r, r, r, x, x,
			},
			secondFrame: []color.Color{
				g, g, g, x, x,
			},
		},
	}

	for _, test := range tests {
		canvas := sparse.NewImage(image.Rect(0, 0, 5, 0))

		stg := New(image.NewRGBA(image.Rect(0, 0, 5, 0)))

		// Draw the first frame of animation.
		drawPixels(test.firstFrame, stg.NextFrame)

		// Commit the first frame.
		stg.Draw(canvas)

		// Draw the second frame.
		drawPixels(test.secondFrame, stg.NextFrame)

		// Commit the second frame.
		stg.Draw(canvas)

		// Confirm the final results match expectations.
		actual := getRow(canvas)

		if !equal(actual, test.secondFrame) {
			t.Errorf("%v: expected %v to be drawn, but got %v", test.name, test.secondFrame, actual)
		}
	}
}

func getRow(img image.Image) []color.Color {
	width := img.Bounds().Dx()
	rv := make([]color.Color, width)
	for x := 0; x < width; x++ {
		rv[x] = img.At(x, 0)
	}
	return rv
}

func equal(actual []color.Color, expected []color.Color) bool {
	for x := 0; x < len(actual); x++ {
		if expected[x] != actual[x] {
			return false
		}
	}
	return true
}

func drawPixels(pixels []color.Color, img draw.Image) {
	for i, p := range pixels {
		img.Set(i, 0, p)
	}
}
