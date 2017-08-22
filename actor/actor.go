package actor

import (
	"image"
	"math"

	"github.com/a-h/raster"
	"github.com/a-h/round"
)

type State struct {
	Angle        float64
	Acceleration float64
	Velocity     float64
	Drag         float64
}

func (s *State) Update(worldBounds image.Rectangle, compositionBounds image.Rectangle, position image.Point, gravity float64) image.Point {
	// Angle remains the same.
	// Acceleration remains the same.
	updatedVelocity := s.Velocity + s.Acceleration - s.Drag
	updatedPosition := calculateNewPosition(position, s.Angle, updatedVelocity)

	// Apply gravity
	updatedPosition.Y += int(gravity)
	// Don't exceed the floor.
	if updatedPosition.Y+compositionBounds.Dy() >= worldBounds.Dy() {
		updatedPosition.Y = worldBounds.Dy() - compositionBounds.Dy()
	}

	s.Velocity = updatedVelocity
	return updatedPosition
}

const (
	degreeToRad = math.Pi / 180
)

func calculateNewPosition(position image.Point, angle float64, velocity float64) image.Point {
	// We know the hypotenuse distance, that the opposite angle is 90 degrees, and
	// we know the angle.
	// soh cah toa

	// To calculate the opposite (y) direction
	// sin(t) = o / h
	// sin(t) * h = o

	// To calculate the adjacent (x) direction
	// cos(t) = a / h
	// cos(t) * h = a
	h := float64(velocity)
	y := math.Sin(degreeToRad*angle) * h
	x := math.Cos(degreeToRad*angle) * h

	toX := position.X + int(round.ToEven(x, 0))
	toY := position.Y + int(round.ToEven(y, 0))

	return image.Point{toX, toY}
}

type CompositionActor struct {
	S *State
	C *raster.Composition
}

func (ca CompositionActor) State() *State {
	return ca.S
}

func (ca CompositionActor) Composition() *raster.Composition {
	return ca.C
}

type Actor interface {
	State() *State
	Composition() *raster.Composition
}
