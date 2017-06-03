package raster

import "image"

func Raycast(current image.Point, p Polygon) int {
	// Keep track of lines we've already hit.
	lines := []*Line{}
	count := 0
	for x := 0; x < current.X; x++ {
		isEdge, isReversal, linesWhichMeet := p.IsEdge(image.Point{x, current.Y})
		if isEdge && !containsAny(lines, linesWhichMeet) {
			count++
		}
		if isReversal {
			count++
		}
		// Keep track of which lines we've already intersected with.
		// Some lines have more than one pixel next to each other.
		if linesWhichMeet != nil && len(linesWhichMeet) > 0 {
			for _, l := range linesWhichMeet {
				lines = append(lines, l)
			}
		}
	}
	return count
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
