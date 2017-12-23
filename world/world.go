package world

import (
	"image"
	"image/draw"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/a-h/raster/actor"
	"github.com/a-h/raster/affine"
)

// World represents a world where things happen. It controls gravity, space and time.
type World struct {
	Background draw.Image
	Actors     []actor.Actor
	Physics    Physics
	Target     draw.Image
	Publisher  Publisher
	Tick       time.Duration
}

type Publisher interface {
	Publish(img draw.Image)
}

type Physics struct {
	Gravity float64
}

func (w *World) Run(stopper <-chan bool) {
	logrus.Debugf("drawing background")
	draw.Draw(w.Target, w.Target.Bounds(), w.Background, image.ZP, draw.Src)
	w.Publisher.Publish(w.Target)

	// Store the locations of each sprite.
	logrus.Debugf("storing sprite locations")
	areas := make([]image.Rectangle, len(w.Actors))
	for i, a := range w.Actors {
		areas[i] = a.Composition().Bounds()
	}

	for {
		select {
		case <-stopper:
			logrus.Debug("received stop signal")
			return
		default:
			logrus.Debugf("drawing background")

			start := time.Now()

			for _, a := range areas {
				// Draw the background over the current location of the actors.
				logrus.Debugf("drawing background over %v", a)
				a = image.Rect(a.Min.X, a.Min.Y, a.Max.X, a.Max.Y)
				draw.Draw(w.Target, a, w.Background, a.Min, draw.Src)
			}

			// Update the display and sleep.
			logrus.Debugf("drawing %d actors", len(w.Actors))
			for i, a := range w.Actors {
				newPosition := a.State().Update(w.Target.Bounds(), a.Composition().Bounds(), a.Composition().Position, w.Physics.Gravity)
				logrus.Debugf("moving actor from %v to %v", a.Composition().Position, newPosition)

				centerX := a.Composition().Bounds().Dx() / 2
				centerY := a.Composition().Bounds().Dy() / 2
				moveToCenter := affine.NewTranslationTransformation(centerX, centerY)
				rotate := affine.NewRotationTransformation(float64(a.State().Angle))
				moveBack := affine.NewTranslationTransformation(-centerX, -centerY)
				a.Composition().Transformation = moveToCenter.Combine(rotate).Combine(moveBack)
				a.Composition().Position = newPosition

				logrus.Debugf("drawing target")
				areas[i] = a.Composition().Draw(w.Target)
				logrus.Debugf("drawn actor at %v", areas[i])
			}

			logrus.Debugf("publishing frame")
			w.Publisher.Publish(w.Target)
			duration := time.Now().Sub(start)
			remaining := w.Tick - duration
			logrus.Debugf("rendered frame in %v of budget %v, %v remaining", duration, w.Tick, remaining)
			time.Sleep(remaining)
		}
	}
}
