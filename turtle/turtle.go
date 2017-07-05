package turtle

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/a-h/raster"
	"github.com/a-h/round"
	"golang.org/x/image/colornames"
)

type Turtle struct {
	Position image.Point
	Angle    float64
	Pen      *Pen
	image    draw.Image
}

type Pen struct {
	Active bool
	Color  color.RGBA
	Size   int
}

func New(image draw.Image) *Turtle {
	return &Turtle{
		image: image,
		Pen: &Pen{
			Active: true,
			Color:  colornames.White,
			Size:   1,
		},
	}
}

const (
	degreeToRad = math.Pi / 180
	radToDegree = 180 / math.Pi
)

func (t *Turtle) Forward(amount int) {
	// We know the hypotenuse distance, that the opposite angle is 90 degrees, and
	// we know the angle.
	// soh cah toa

	// To calculate the opposite (y) direction
	// sin(t) = o / h
	// sin(t) * h = o

	// To calculate the adjacent (x) direction
	// cos(t) = a / h
	// cos(t) * h = a
	h := float64(amount)
	y := math.Sin(degreeToRad*t.Angle) * h
	x := math.Cos(degreeToRad*t.Angle) * h

	toX := t.Position.X + int(round.ToEven(x, 0))
	toY := t.Position.Y + int(round.ToEven(y, 0))

	to := image.Point{toX, toY}

	if t.Pen.Active {
		l := raster.NewLine(t.Position, to, t.Pen.Color)
		l.Draw(t.image)
	}

	t.Position = to
}

func (t *Turtle) Rotate(degrees float64) {
	n := t.Angle + degrees
	if n > 360 {
		n -= 360
	}
	if n < 0 {
		n += 360
	}
	t.Angle = n
}
