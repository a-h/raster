package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"github.com/a-h/raster"
	"golang.org/x/image/colornames"

	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/cocoa"
)

const (
	height = 600
	width  = 800
)

func main() {
	go run()
	wde.Run()
}

func run() {
	w, err := wde.NewWindow(width, height)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("failed to create window: %v", err))
		os.Exit(-1)
	}
	defer w.Close()
	defer wde.Stop()
	w.SetSize(width, height)
	w.Show()

	img := w.Screen()

	wmiddle := width / 2
	hmiddle := height / 2
	wpinch := wmiddle / 4
	hpinch := hmiddle / 4
	a := image.Point{wmiddle, 0}
	b := image.Point{wmiddle + wpinch, hmiddle - hpinch}
	c := image.Point{width, hmiddle}
	d := image.Point{wmiddle + wpinch, hmiddle + hpinch}
	e := image.Point{wmiddle, height}
	f := image.Point{wmiddle - wpinch, hmiddle + hpinch}
	g := image.Point{0, hmiddle}
	h := image.Point{wmiddle - wpinch, hmiddle - hpinch}

	colors := []color.RGBA{colornames.Gray, colornames.Red, colornames.Green, colornames.Blue}
	i := 0
	for {
		color := colors[i%len(colors)]
		p := raster.NewPolygon(color, a, b, c, d, e, f, g, h)
		p.Draw(img)
		w.FlushImage()
		i++
		time.Sleep(100)
	}
}
