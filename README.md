## Raster

A very simple library for drawing 2D shapes onto images.

I wrote it to learn, and to render simple shapes for games, you're probably looking for something like this: https://github.com/fogleman/gg

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

## Examples

See [./examples](examples).