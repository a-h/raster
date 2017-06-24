package biggest

import "testing"

func TestIntegerFunction(t *testing.T) {
	tests := []struct {
		input    []int
		expected int
	}{
		{
			input:    []int{1, 2, 3},
			expected: 3,
		},
		{
			input:    []int{3, 2, 1},
			expected: 3,
		},
	}

	for _, test := range tests {
		actual := IntegerIn(test.input...)
		if actual != test.expected {
			t.Errorf("for input %v, expected %v, got %v", test.input, test.expected, actual)
		}
	}
}
