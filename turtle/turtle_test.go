package turtle

import (
	"image"
	"testing"
)

func TestThatPositionsAreCalculated(t *testing.T) {
	tests := []struct {
		Angle    float64
		Forward  int
		Expected image.Point
	}{
		{
			Angle:    0,
			Forward:  5,
			Expected: image.Point{5, 0},
		},
		{
			Angle:    0,
			Forward:  -5,
			Expected: image.Point{-5, 0},
		},
		{
			Angle:    90,
			Forward:  5,
			Expected: image.Point{0, 5},
		},
		{
			Angle:    180,
			Forward:  5,
			Expected: image.Point{-5, 0},
		},
		{
			Angle:    270,
			Forward:  5,
			Expected: image.Point{0, -5},
		},
		{
			Angle:    90 / 2, // Equal sided, right angle
			Forward:  5,
			Expected: image.Point{4, 4}, // Actual value is 3.5, it's rounded to 4.
		},
		{
			Angle:    36.87,
			Forward:  5,
			Expected: image.Point{4, 3},
		},
		{
			Angle:    53.130102,
			Forward:  5,
			Expected: image.Point{3, 4},
		},
	}

	for _, test := range tests {
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		o := New(img)
		o.Angle = test.Angle
		o.Forward(test.Forward)
		actual := o.Position

		if !actual.Eq(test.Expected) {
			t.Errorf("for %v∘, moving forward %v pixels. expected %v, got %v", test.Angle, test.Forward, test.Expected, actual)
		}
	}
}

func TestThatTheTurtleCanRotate(t *testing.T) {
	tests := []struct {
		Angle    float64
		Rotation float64
		Forward  int
		Expected image.Point
	}{
		{
			Angle:    270,
			Rotation: 90,
			Forward:  5,
			Expected: image.Point{5, 0},
		},
		{
			Angle:    270,
			Rotation: 90,
			Forward:  -5,
			Expected: image.Point{-5, 0},
		},
		{
			Angle:    0,
			Rotation: 90,
			Forward:  5,
			Expected: image.Point{0, 5},
		},
		{
			Angle:    90,
			Rotation: 90,
			Forward:  5,
			Expected: image.Point{-5, 0},
		},
		{
			Angle:    180,
			Rotation: 90,
			Forward:  5,
			Expected: image.Point{0, -5},
		},
		{
			Angle:    0,
			Rotation: 36.87,
			Forward:  5,
			Expected: image.Point{4, 3},
		},
		{
			Angle:    0,
			Rotation: 53.130102,
			Forward:  5,
			Expected: image.Point{3, 4},
		},
		{
			Angle:    0,
			Rotation: 90 + 360,
			Forward:  5,
			Expected: image.Point{0, 5},
		},
		{
			Angle:    0,
			Rotation: 90 - 360,
			Forward:  5,
			Expected: image.Point{0, 5},
		},
	}

	for _, test := range tests {
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		o := New(img)
		o.Angle = test.Angle
		o.Rotate(test.Rotation)
		o.Forward(test.Forward)
		actual := o.Position

		if !actual.Eq(test.Expected) {
			t.Errorf("starting at %v∘, then rotating by %v∘, then moving forward by %v pixels. expected %v, got %v", test.Angle, test.Rotation, test.Forward, test.Expected, actual)
		}
	}
}

func TestThatTheTurtleCanMoveBackwards(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	o := New(img)
	o.Forward(10)
	if !o.Position.Eq(image.Point{10, 0}) {
		t.Error("expected to be able to move forward")
	}
	o.Backward(5)
	if !o.Position.Eq(image.Point{5, 0}) {
		t.Error("expected to be able to move backward")
	}
	o.Forward(5)
	if !o.Position.Eq(image.Point{10, 0}) {
		t.Error("expected angle to be preserved")
	}
	o.Backward(10)
	if !o.Position.Eq(image.Point{0, 0}) {
		t.Error("expected to be able to move back to the origin")
	}
}
