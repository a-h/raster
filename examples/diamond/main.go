package main

import (
	"fmt"
	"image"

	"github.com/a-h/raster"
	"github.com/a-h/raster/affine"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/colornames"
	"golang.org/x/mobile/event/lifecycle"
)

const (
	height = 500
	width  = 500
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

		composition := raster.NewComposition(image.Point{0, 0},
			raster.NewPolygon(colornames.White, image.Point{250, 0}, image.Point{500, 250}, image.Point{250, 500}, image.Point{0, 250}))

		composition.Transformation = affine.NewRotationTransformation(45)

		composition.Draw(img)

		w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, width, height))
		w.Publish()

		// Keep looking for events.
		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					fmt.Println("StageDead...")
					return
				}
			case error:
				fmt.Println("Error: ", e.Error())
			}
		}
	})
}
