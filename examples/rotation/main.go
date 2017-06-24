package main

import (
	"fmt"
	"image"

	"github.com/a-h/raster"
	"github.com/a-h/raster/affine"

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

		l := raster.NewLine(image.Point{0, 0}, image.Point{250, 0}, colornames.Green)

		// Draw a plus symbol.
		c := raster.NewComposition(image.Point{500, 500}, l)
		c.Draw(img)
		c.Transformation = affine.NewRotationTransformation(90)
		c.Draw(img)
		c.Transformation = affine.NewRotationTransformation(180)
		c.Draw(img)
		c.Transformation = affine.NewRotationTransformation(270)
		c.Draw(img)

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
