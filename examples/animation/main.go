package main

import (
	"fmt"
	"image"
	"image/draw"

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

		for i := 0; i <= 1000; i += 5 {
			// Draw the items on the next frame of the stage.
			square := raster.NewSquare(image.Point{i, 250}, 500, colornames.Green)
			square.Draw(stg.NextFrame)

			circle := raster.NewCircle(image.Point{i + 50, i + 50}, 50, colornames.Red)
			circle.Draw(stg.NextFrame)

			// Draw the whole stage.
			stg.Draw(img)
			w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
			w.Publish()
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
