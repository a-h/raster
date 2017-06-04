package raster

import "image"

func ScanLine(y int, p Polygon) []int {
	// Keep track of lines we've already hit.
	width := p.Bounds().Dx()
	rv := make([]int, width)
	lines := []*Line{}
	count := 0

	for x := 0; x < p.Bounds().Dx(); x++ {
		isEdge, linesWhichMeet := p.IsEdge(image.Point{x, y})
		if isEdge && !containsAny(lines, linesWhichMeet) {
			count++
		}
		// Keep track of which lines we've already intersected with.
		// Some lines have more than one pixel next to each other.
		if linesWhichMeet != nil && len(linesWhichMeet) > 0 {
			for _, l := range linesWhichMeet {
				lines = append(lines, l)
			}
		}

		if len(linesWhichMeet) > 1 {
			for i := 0; i < len(linesWhichMeet); i += 2 {
				l1 := linesWhichMeet[i]
				l2 := linesWhichMeet[i+1]
				d := CalculateDirection(l1, l2)
				if d == Up || d == Down {
					// Count up and down vertices twice.
					count++
				}
			}
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
