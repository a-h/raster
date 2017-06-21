package raster

import "testing"

func TestThatCirclesAreComposable(t *testing.T) {
	var c interface{} = new(Circle)
	if _, ok := c.(Composable); !ok {
		t.Errorf("expected %t to implement Composable", c)
	}
}

func TestThatCompositionsAreComposable(t *testing.T) {
	var c interface{} = new(Composition)
	if _, ok := c.(Composable); !ok {
		t.Errorf("expected %t to implement Composable", c)
	}
}

func TestThatLinesAreComposable(t *testing.T) {
	var c interface{} = new(Line)
	if _, ok := c.(Composable); !ok {
		t.Errorf("expected %t to implement Composable", c)
	}
}

func TestThatPolygonsAreComposable(t *testing.T) {
	var c interface{} = new(Polygon)
	if _, ok := c.(Composable); !ok {
		t.Errorf("expected %t to implement Composable", c)
	}
}

func TestThatSquaresAreComposable(t *testing.T) {
	var c interface{} = new(Square)
	if _, ok := c.(Composable); !ok {
		t.Errorf("expected %t to implement Composable", c)
	}
}

func TestThatTextIsComposable(t *testing.T) {
	var c interface{} = new(Text)
	if _, ok := c.(Composable); !ok {
		t.Errorf("expected %t to implement Composable", c)
	}
}
