package raster

import "testing"

func TestThatCirclesAreComposable(t *testing.T) {
	var c interface{} = new(Circle)
	if _, ok := c.(Composable); !ok {
		t.Error("expected Circle to implement Composable")
	}
}

func TestThatFilledCirclesAreComposable(t *testing.T) {
	var c interface{} = new(FilledCircle)
	if _, ok := c.(Composable); !ok {
		t.Error("expected Filled Circle to implement Composable")
	}
}

func TestThatCompositionsAreComposable(t *testing.T) {
	var c interface{} = new(Composition)
	if _, ok := c.(Composable); !ok {
		t.Error("expected Composition to implement Composable")
	}
}

func TestThatLinesAreComposable(t *testing.T) {
	var c interface{} = new(Line)
	if _, ok := c.(Composable); !ok {
		t.Error("expected Lines to implement Composable")
	}
}

func TestThatPolygonsAreComposable(t *testing.T) {
	var c interface{} = new(Polygon)
	if _, ok := c.(Composable); !ok {
		t.Error("expected Polygons to implement Composable")
	}
}

func TestThatFilledPolygonsAreComposable(t *testing.T) {
	var c interface{} = new(FilledPolygon)
	if _, ok := c.(Composable); !ok {
		t.Error("expected Filled Polygons to implement Composable")
	}
}

func TestThatSquaresAreComposable(t *testing.T) {
	var c interface{} = new(Square)
	if _, ok := c.(Composable); !ok {
		t.Error("expected Squares to implement Composable")
	}
}

func TestThatTextIsComposable(t *testing.T) {
	var c interface{} = new(Text)
	if _, ok := c.(Composable); !ok {
		t.Error("expected Text to implement Composable")
	}
}
