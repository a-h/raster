package actor

import (
	"image"
	"math"

	"github.com/a-h/raster"
	"github.com/a-h/round"
)

// State represents the current state of an actor, it's direction (angle), acceleration, velocity, and drag.
type State struct {
	Angle        float64
	Acceleration float64
	Velocity     float64
	Drag         float64
}

// Update updates the state's velocity and calculates a new position for the object.
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

// CompositionActor represents a Composition with state. It's the simplest implementation of
// the Actor interface, which defines State (direction, velocity etc.) and the Composition which
// defines the pixels to display.
type CompositionActor struct {
	S *State
	C *raster.Composition
}

// State represents the current state of an actor, it's direction (angle), acceleration, velocity, and drag.
func (ca CompositionActor) State() *State {
	return ca.S
}

// Composition returns the position and components which make it up, and a transformation
// that can be applied to it to move, scale, or rotate all of the elements.
func (ca CompositionActor) Composition() *raster.Composition {
	return ca.C
}

// Actor defines an item to place on a stage. It has a State (direction, velocity) and
// composition which defines the pixels drawn.
type Actor interface {
	State() *State
	Composition() *raster.Composition
}
