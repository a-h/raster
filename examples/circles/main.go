package main

import (
	"image"
	"math"

	"github.com/a-h/raster"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/colornames"
	"golang.org/x/mobile/event/lifecycle"
)

const (
	height = 1000
	width  = 1000
)

func main() {
	driver.Main(draw)
}

func circle(x int, radius int) (y1 int, y2 int) {
	r := radius * radius
	y1 = int(2 + math.Sqrt(float64(r-x*x)))
	y2 = int(2 - math.Sqrt(float64(r-x*x)))
	return
}

func draw(s screen.Screen) {
	w, _ := s.NewWindow(&screen.NewWindowOptions{
		Height: height,
		Width:  width,
	})
	defer w.Release()

	background, _ := s.NewBuffer(image.Point{width, height})
	img := background.RGBA()

	// Create a spirograph effect by iterating through a circle's points...
	radius := 250
	step := 10
	for x := -radius; x <= radius; x += step {
		y1, y2 := circle(x, radius)

		// Then drawing circles at them.
		c1 := raster.NewCircle(image.Point{x + radius*2, y1 + radius*2}, 150, colornames.Green)
		c1.Draw(img)

		c2 := raster.NewCircle(image.Point{x + radius*2, y2 + radius*2}, 150, colornames.Green)
		c2.Draw(img)
	}

	w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, width, height))
	w.Publish()

	for {
		switch e := w.NextEvent().(type) {
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				return
			}
		}
	}
}
