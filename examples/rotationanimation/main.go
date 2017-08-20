package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/a-h/raster/affine"
	"github.com/a-h/raster/stage"

	"github.com/a-h/raster"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/colornames"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
)

const (
	height = 1000
	width  = 1000
)

func main() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Height: height,
			Width:  width,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		defer w.Release()

		background, _ := s.NewBuffer(image.Point{width, height})
		img := background.RGBA()

		stgBackdrop := image.NewRGBA(img.Bounds())
		// Set the background as white.
		draw.Draw(stgBackdrop, stgBackdrop.Bounds(), &image.Uniform{colornames.White}, image.ZP, draw.Src)
		stg := stage.New(stgBackdrop)

		// Use a composition to hold a triangle and a square to show the bounds.
		triangle := raster.NewPolygon(colornames.Black, image.Point{0, 100}, image.Point{50, 0}, image.Point{100, 100})
		square := raster.NewPolygon(colornames.Lightgrey, image.Point{0, 0}, image.Point{100, 0}, image.Point{100, 100}, image.Point{0, 100})
		c := raster.NewComposition(image.Point{500 - 50, 500 - 50}, triangle, square)

		for i := 0; i <= 1000; i += 5 {
			// Draw the items on the next frame of the stage.
			c.Draw(stg.NextFrame)

			// Draw the whole stage.
			stg.Draw(img)
			w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
			w.Publish()

			// Rotate the triangle from the center.
			moveToCenter := affine.NewTranslationTransformation(50, 50)
			rotate := affine.NewRotationTransformation(float64(i))
			moveBack := affine.NewTranslationTransformation(-50, -50)
			c.Transformation = moveToCenter.Combine(rotate).Combine(moveBack)
		}

		// Keep looking for events.
		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					fmt.Println("StageDead...")
					return
				}
			case key.Event:
				if e.Code == key.CodeUpArrow {
					// Do something.
				}
			case error:
				fmt.Println("Error: ", e.Error())
			}
		}
	})
}
