package raster

import "image"

func ScanLine(y int, p Polygon) []int {
	// Keep track of lines we've already hit.
	width := p.Bounds().Dx()
	rv := make([]int, width)
	lines := []*Line{}
	count := 0

	plines := map[*Line]int{}
	for i, l := range p.Lines {
		plines[l] = (i % 2) + 1
	}

	for x := 0; x < p.Bounds().Dx(); x++ {
		isEdge, linesWhichMeet := p.IsEdge(image.Point{x, y})

		for _, l := range linesWhichMeet {
			if !contains(lines, l) {
				count += plines[l]
			}

			lines = append(lines, l)
		}

		if isEdge {
			rv[x] = 0 // Don't overwrite the edges
		} else {
			rv[x] = count
		}
	}
	return rv
}

// Direction represents the direction a vertex (point) makes when two lines come together.
type Direction int

const (
	// None is no direction
	None = Direction(0)
	// Up points North
	Up = Direction(1)
	// Right points East
	Right = Direction(2)
	// Down points South
	Down = Direction(3)
	// Left points West
	Left = Direction(4)
)

func CalculateDirection(l1 *Line, l2 *Line) Direction {
	var sharedPoint = l1.To
	if !sharedPoint.Eq(l2.From) {
		return None
	}
	sp := l1.From
	ep := l2.To
	// If the X point of both is to the right of the shared point, it's a left arrow
	if sp.X > sharedPoint.X && ep.X > sharedPoint.X {
		return Left
	}
	// If the X point of both is to the left, it's a right arrow
	if sp.X < sharedPoint.X && ep.X < sharedPoint.X {
		return Right
	}
	// If the end Y point of both lines is to the bottom of the shared point, it's an up arrow
	if sp.Y > sharedPoint.Y && ep.Y > sharedPoint.Y {
		return Up
	}
	// If the end Y point of both lines is to the top of the shared point, it's a down arrow
	if sp.Y < sharedPoint.Y && ep.Y < sharedPoint.Y {
		return Down
	}
	return None
}

func containsAny(lines []*Line, of []*Line) bool {
	for _, ll := range lines {
		for _, jj := range of {
			if ll.Eq(jj) {
				return true
			}
		}
	}
	return false
}

func contains(lines []*Line, l *Line) bool {
	for _, ll := range lines {
		if ll.Eq(l) {
			return true
		}
	}
	return false
}
