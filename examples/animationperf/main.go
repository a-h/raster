package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/colornames"
)

func main() {
	// This example shows that the screen driver (on OSX at least)
	// takes longer to upload and publish the more changes are made
	// per frame.

	fmt.Println("Testing redrawing 1000x1000 pixels each frame.")
	test(true, 1000)

	/*
		fmt.Println("Testing only changing a single pixel each frame.")
		test(false, 1000)
	*/
}

func test(redraw bool, iterations int) {
	width, height := 1000, 1000
	uploadTimes := []time.Duration{}
	publishTimes := []time.Duration{}

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

		for i := 0; i < iterations; i++ {
			if redraw {
				// Draw over the whole image.
				for y := 0; y <= height; y++ {
					for x := 0; x <= width; x++ {
						img.Set(x, y, randomColor())
					}
				}
			} else {
				// Just draw a single pixel.
				x, y := rand.Intn(width), rand.Intn(height)
				img.Set(x, y, randomColor())
			}

			uploadStart := time.Now()
			w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, width, height))
			if i > 0 {
				uploadTimes = append(uploadTimes, time.Since(uploadStart))
			}

			publishStart := time.Now()
			w.Publish()
			if i > 0 {
				publishTimes = append(publishTimes, time.Since(publishStart))
			}
		}

		fmt.Printf("Upload: %v, Publish: %v\n", sum(uploadTimes), sum(publishTimes))
	})
}

func average(durations []time.Duration) time.Duration {
	l := len(durations)
	s := sum(durations)

	return time.Duration(int64(s) / int64(l))
}

func sum(durations []time.Duration) time.Duration {
	var r time.Duration

	for i := 0; i < len(durations); i++ {
		r += durations[i]
	}

	return r
}

func randomColor() color.RGBA {
	switch rand.Intn(5) {
	case 0:
		return colornames.Red
	case 1:
		return colornames.Green
	case 2:
		return colornames.Blue
	case 3:
		return colornames.Yellow
	case 4:
		return colornames.Cyan
	case 5:
		return colornames.Purple
	}
	return colornames.White
}
