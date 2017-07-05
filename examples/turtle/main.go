package main

import (
	"fmt"
	"image"

	"github.com/a-h/raster/turtle"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
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

		t := turtle.New(img)
		t.Position = image.Point{width / 2, (height / 2) - 60}
		// Walk out a bit.
		t.Pen.Active = false
		t.Forward(30)
		// Then draw an octogon.
		sides := 8
		angle := 360 / float64(sides)

		t.Pen.Active = true
		for i := 0; i < sides; i++ {
			t.Rotate(angle)
			t.Forward(50)
		}

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
