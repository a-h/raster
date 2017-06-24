package smallest

// Integer returns the smallest int in the array.
func Integer(ints []int) int {
	s := ints[0]
	for _, v := range ints[1:] {
		if v < s {
			s = v
		}
	}
	return s
}

// IntegerIn returns the smallest integer from the inputs.
func IntegerIn(ints ...int) int {
	return Integer(ints)
}
