package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/skelterjohn/go.wde"

	// OSX only, need to use a different backend for other operating systems.
	_ "github.com/skelterjohn/go.wde/cocoa"

	"github.com/a-h/raster/actor"
	"github.com/a-h/raster/world"

	"github.com/a-h/raster"

	"golang.org/x/image/colornames"
)

const (
	height = 500
	width  = 500
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	go run()
	wde.Run()
}

func run() {
	// Create a WDE window.
	logrus.Debugf("creating new window")
	w, err := wde.NewWindow(width, height)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("failed to create window: %v", err))
		os.Exit(-1)
	}
	defer w.Close()
	w.SetSize(width, height)
	w.Show()

	// The buffer to write to.
	img := w.Screen()

	// Create a background.
	logrus.Debugf("creating background")
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(background, background.Bounds(), image.NewUniform(colornames.White), image.ZP, draw.Src)

	// Create a composition to use.
	logrus.Debugf("creating composition")
	c := raster.NewCircle(image.Point{0, 50}, 50, colornames.Black)
	ball := &actor.CompositionActor{
		C: raster.NewComposition(image.Point{250, 0}, c),
		S: &actor.State{
			Angle:        0.0,
			Velocity:     0.0,
			Acceleration: 0.0,
			Drag:         0.0,
		},
	}

	logrus.Debugf("creating publisher")
	p := NewWdePublisher(w)

	nw := world.World{
		Background: background,
		Actors:     []actor.Actor{ball},
		Physics: world.Physics{
			Gravity: 5,
		},
		Target:    img,
		Publisher: p,
	}

	logrus.Debugf("starting world")
	stopper := make(<-chan bool)
	nw.Run(stopper)
	wde.Stop()
}

type WdePublisher struct {
	Window wde.Window
}

func (w *WdePublisher) Publish(img draw.Image) {
	w.Window.FlushImage()
}

func NewWdePublisher(w wde.Window) *WdePublisher {
	return &WdePublisher{
		Window: w,
	}
}
