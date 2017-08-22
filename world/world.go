package world

import (
	"image"

	"github.com/Sirupsen/logrus"
	"golang.org/x/exp/shiny/screen"

	"github.com/a-h/raster/actor"
	"github.com/a-h/raster/affine"
	"github.com/a-h/raster/stage"
)

// World represents a world where things happen. It controls gravity, space and time.
type World struct {
	Actors  []actor.Actor
	Stage   *stage.Stage
	Physics Physics
	Buffer  screen.Buffer
	Window  screen.Window
}

type Physics struct {
	Gravity float64
}

func (w *World) Run(stopper <-chan bool) {
	bufferImage := w.Buffer.RGBA()

	for {
		select {
		case <-stopper:
			logrus.Debug("received stop signal")
			return
		default:
			// Update the display and sleep.
			for _, a := range w.Actors {
				newPosition := a.State().Update(w.Stage.Bounds(), a.Composition().Bounds(), a.Composition().Position, w.Physics.Gravity)

				centerX := a.Composition().Bounds().Dx() / 2
				centerY := a.Composition().Bounds().Dy() / 2
				moveToCenter := affine.NewTranslationTransformation(centerX, centerY)
				rotate := affine.NewRotationTransformation(float64(a.State().Angle))
				moveBack := affine.NewTranslationTransformation(-centerX, -centerY)
				a.Composition().Transformation = moveToCenter.Combine(rotate).Combine(moveBack)
				a.Composition().Position = newPosition

				a.Composition().Draw(w.Stage.NextFrame)
			}

			w.Stage.Draw(bufferImage)
			w.Window.Upload(image.Point{0, 0}, w.Buffer, bufferImage.Bounds())
			w.Window.Publish()
			logrus.Debugf("Rendered frame.\n")
		}
	}
}
