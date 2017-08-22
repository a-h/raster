package stage

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/a-h/raster/sparse"
)

// Stage is a stage for creating animation frames on top of.
type Stage struct {
	Backdrop     image.Image
	CurrentFrame *sparse.Image
	NextFrame    *sparse.Image
	// Set to true when the first frame has been drawn.
	frameDrawn bool
}

// New creates a new Stage, initialised with a backdrop.
func New(backdrop image.Image) *Stage {
	return &Stage{
		Backdrop:     backdrop,
		CurrentFrame: sparse.NewImage(backdrop.Bounds()),
		NextFrame:    sparse.NewImage(backdrop.Bounds()),
	}
}

// ResetCurrentFrame resets the frame that's currently on screen.
func (stg *Stage) ResetCurrentFrame() {
	stg.CurrentFrame = sparse.NewImage(stg.Backdrop.Bounds())
}

// ResetNextFrame clears the frame that's going to be drawn next.
func (stg *Stage) ResetNextFrame() {
	stg.NextFrame = sparse.NewImage(stg.Backdrop.Bounds())
}

// Draw unpaints draws pixels in the ToDraw field onto the Backdrop and adds them to the Drawn field.
// At the end, the ToDraw field is wiped. It's assumed that nothing else alters the img parameter other
// than the stage itself.
func (stg *Stage) Draw(img draw.Image) (redraw bool) {
	// If it's the first draw, draw the entire backdrop onto the image.
	if !stg.frameDrawn {
		b := stg.Backdrop.Bounds()
		for y := 0; y < b.Dy(); y++ {
			for x := 0; x < b.Dx(); x++ {
				img.Set(x, y, stg.Backdrop.At(x, y))
			}
		}
		stg.frameDrawn = true
		// Tell the caller to redraw the whole frame.
		redraw = true
	}

	// Set everything that's already been drawn back to the background color, but only
	// only when the NextFrame isn't going to set it to another (or the same) color.
	pixelsToSkip := make(map[image.Point]interface{})
	for p, existingColor := range stg.CurrentFrame.Drawn {
		if newColor, pixelWillBeRedrawn := stg.NextFrame.Drawn[p]; pixelWillBeRedrawn {
			// The next pixel is the same color anyway, so skip it.
			if newColor == existingColor {
				pixelsToSkip[p] = true
			}
		} else {
			// Set the pixel back to the backdrop color.
			backgroundColor := stg.Backdrop.At(p.X, p.Y)
			img.Set(p.X, p.Y, backgroundColor)
		}
	}

	// Draw the pixels that need it back out.
	for p, c := range stg.NextFrame.Drawn {
		// If it's already been drawn, or the color is the same as the backdrop, don't ask
		// the image to set it again.
		if _, shouldSkip := pixelsToSkip[p]; !shouldSkip {
			img.Set(p.X, p.Y, c)
		}
	}

	// Copy pixels from the next frame to the current frame to keep a representation
	// of what should be on the image.
	stg.ResetCurrentFrame()
	for p, c := range stg.NextFrame.Drawn {
		stg.CurrentFrame.Drawn[p] = c
	}
	stg.ResetNextFrame()
	return
}

// Implementation of the image.Draw interface below.
// Stage draws the "CurrentFrame".

// ColorModel returns the Image's color model.
func (stg *Stage) ColorModel() color.Model {
	return stg.Backdrop.ColorModel()
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (stg *Stage) Bounds() image.Rectangle {
	return stg.Backdrop.Bounds()
}

// Set the pixel with x, y coordinates to the color c.
// The results of this are essentially thrown away.
func (stg *Stage) Set(x, y int, c color.Color) {
	stg.CurrentFrame.Drawn[image.Point{x, y}] = c
}

// At returns the color of the pixel at (x, y).
func (stg *Stage) At(x, y int) color.Color {
	p := image.Point{x, y}
	if c, ok := stg.CurrentFrame.Drawn[p]; ok {
		return c
	}
	return stg.Backdrop.At(x, y)
}
