package main

import (
	"fmt"
	"image"

	"github.com/a-h/raster"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/colornames"
	"golang.org/x/mobile/event/lifecycle"
)

const (
	height      = 1000
	width       = 1000
	squareWidth = 100
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

		// Set the background colour.
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				img.Set(x, y, colornames.Lightgoldenrodyellow)
			}
		}

		angle := 10

		for i := 0; i < width-squareWidth; i += squareWidth {
			for j := 0; j < height-squareWidth; j += squareWidth {
				ax, ay := i, j
				bx, by := i+squareWidth-angle, j+angle
				cx, cy := i+squareWidth, j+squareWidth
				dx, dy := i+angle, j+squareWidth-angle

				a := image.Point{int(ax), int(ay)}
				b := image.Point{int(bx), int(by)}
				c := image.Point{int(cx), int(cy)}
				d := image.Point{int(dx), int(dy)}

				p := raster.NewFilledPolygon(colornames.Orange, colornames.Darkorange, a, b, c, d)
				p.Draw(img)

				angle++
			}
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
