package smallest

import "testing"

func TestIntegerFunction(t *testing.T) {
	tests := []struct {
		input    []int
		expected int
	}{
		{
			input:    []int{-1, 2, 3},
			expected: -1,
		},
		{
			input:    []int{3, 2, 1},
			expected: 1,
		},
	}

	for _, test := range tests {
		actual := IntegerIn(test.input...)
		if actual != test.expected {
			t.Errorf("for input %v, expected %v, got %v", test.input, test.expected, actual)
		}
	}
}
