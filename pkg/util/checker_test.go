package util

import "testing"

func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{nil, true},
		{"", true},
		{"   ", true},
		{"non-empty", false},
		{[]int{}, true},
		{[]int{1, 2, 3}, false},
		{map[string]int{}, true},
		{map[string]int{"key": 1}, false},
		{0, true},
		{1, false},
		{0.0, true},
		{1.5, false},
		{false, true},
		{true, false},
		{struct{}{}, false},
		{(*int)(nil), true},
	}

	for _, tt := range tests {
		t.Run("IsEmptyValue", func(t *testing.T) {
			result := IsEmptyValue(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for input %v", tt.expected, result, tt.input)
			}
		})
	}
}
