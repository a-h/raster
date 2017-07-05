## Raster

A very simple library for drawing 2D shapes onto images or the screen with minimal dependencies.

I wrote it to learn, and to render simple shapes for games. You're probably looking for something more fully featured, like this: https://github.com/fogleman/gg

There's also a Turtle library similar to PyTurtle I wrote for my son to play with.

### Features

```go
// Circle
circle := raster.NewCircle(image.Point{500, 500}, 250, colornames.Maroon)
circle.Draw(img)

// Square
square := raster.NewSquare(image.Point{250, 250}, 500, colornames.Green)
square.Draw(img)

// Polygon
a := image.Point{500, 0}
b := image.Point{750, 250}
c := image.Point{1000, 500}
d := image.Point{750, 750}
e := image.Point{500, 1000}
f := image.Point{250, 750}
g := image.Point{0, 500}
h := image.Point{250, 250}

p := raster.NewFilledPolygon(colornames.Gray, colornames.Antiquewhite, a, b, c, d, e, f, g, h)
p.Draw(img)

// Text
t := raster.NewText(image.Point{0, 0}, "Hello!", colornames.White)
t.Draw(img)

// Combine elements together.
circleInsideSquare := raster.NewComposition(image.Point{250, 250},
    raster.NewCircle(image.Point{250, 250}, 250, colornames.Maroon),
    raster.NewSquare(image.Point{0, 0}, 500, colornames.Green))
circleInsideSquare.Draw(img)
```

### Turtle

See [./examples/turtle](./examples/turtle)

```go
t := turtle.New(img)
t.Position = image.Point{width / 2, (height / 2) - 60}
// Walk out a bit.
t.Pen.Active = false
t.Forward(30)
// Then draw an octogon.
sides := 8
angle := 360 / float64(sides)

t.Pen.Active = true
for i := 0; i < sides; i++ {
    t.Rotate(angle)
    t.Forward(50)
}
```


## Complete examples

See [./examples](examples)