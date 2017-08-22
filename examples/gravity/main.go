package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/a-h/raster/actor"
	"github.com/a-h/raster/stage"
	"github.com/a-h/raster/world"

	"github.com/a-h/raster"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/colornames"
)

const (
	height = 800
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

		// Create a composition to use.
		c := raster.NewCircle(image.Point{0, 50}, 50, colornames.Black)
		ball := &actor.CompositionActor{
			C: raster.NewComposition(image.Point{500, 0}, c),
			S: &actor.State{
				Angle:        0.0,
				Velocity:     0.0,
				Acceleration: 0.0,
				Drag:         0.0,
			},
		}

		nw := world.World{
			Actors: []actor.Actor{ball},
			Stage:  stg,
			Physics: world.Physics{
				Gravity: 5,
			},
			Window: w,
			Buffer: background,
		}

		stopper := make(<-chan bool)
		nw.Run(stopper)
	})
}
