package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/a-h/raster"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/colornames"
	"golang.org/x/mobile/event/key"
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

		angle := 10

		for i := 0; i < width-squareWidth; i += squareWidth {
			for j := 0; j < height-squareWidth; j += squareWidth {
				ax, ay := i, j // 0, 0
				bx, by := i+squareWidth-angle, j+angle
				cx, cy := i+squareWidth, j+squareWidth
				dx, dy := i+angle, j+squareWidth-angle

				for _, p := range raster.Line(ax, ay, bx, by) {
					img.SetRGBA(p.X, p.Y, colornames.Red)
				}

				for _, p := range raster.Line(bx, by, cx, cy) {
					img.SetRGBA(p.X, p.Y, colornames.Green)
				}

				for _, p := range raster.Line(cx, cy, dx, dy) {
					img.SetRGBA(p.X, p.Y, colornames.Lightblue)
				}

				for _, p := range raster.Line(dx, dy, ax, ay) {
					img.SetRGBA(p.X, p.Y, colornames.Yellow)
				}

				/*
					a := image.Point{int(ax), int(ay)}
					b := image.Point{int(bx), int(by)}
					c := image.Point{int(cx), int(cy)}
					d := image.Point{int(dx), int(dy)}

					raster.DrawFilledPolygon(img, colornames.Red, colornames.Blue, a, b, c, d)
				*/

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
			case key.Event:
				if e.Code == key.CodeUpArrow {
					// Do something.
				}
			case error:
			}
		}
	})
}

func drawBackground(img *image.RGBA, c color.RGBA) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, c)
		}
	}
}
