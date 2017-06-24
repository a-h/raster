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

		middle := width / 2
		pinch := middle - 400
		a := image.Point{middle, 0}
		b := image.Point{middle + pinch, middle - pinch}
		c := image.Point{1000, 500}
		d := image.Point{middle + pinch, middle + pinch}
		e := image.Point{middle, 1000}
		f := image.Point{middle - pinch, middle + pinch}
		g := image.Point{0, 500}
		h := image.Point{middle - pinch, middle - pinch}

		/*
			c := image.Point{1000, 500}
			d := image.Point{750, 750}
			e := image.Point{500, 1000}
			f := image.Point{250, 750}
			g := image.Point{0, 500}
			h := image.Point{250, 250}
		*/

		p := raster.NewFilledPolygon(colornames.Gray, colornames.Antiquewhite, a, b, c, d, e, f, g, h)
		p.Draw(img)

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
