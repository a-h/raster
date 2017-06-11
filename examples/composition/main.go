package main

import (
	"fmt"
	"image"

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

		circleInsideSquare := raster.NewComposition(image.Point{250, 250},
			raster.NewCircle(250, 250, 250, colornames.Maroon),
			raster.NewSquare(0, 0, 500, colornames.Green))
		circleInsideSquare.Draw(img)

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
				if e.To == lifecycle.StageVisible {
					fmt.Println("StageVisible...")
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
