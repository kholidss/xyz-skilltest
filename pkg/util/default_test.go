package util

import "testing"

func TestSetDefaultInt(t *testing.T) {
	tests := []struct {
		value        any
		defaultValue int
		expected     int
	}{
		{nil, 10, 10},
		{"string", 10, 10},
		{0, 10, 10},
		{5, 10, 5},
		{42, 10, 42},
		{3.14, 10, 3},
	}

	for _, tt := range tests {
		t.Run("SetDefaultInt", func(t *testing.T) {
			result := SetDefaultInt(tt.value, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %d, got %d for input %v", tt.expected, result, tt.value)
			}
		})
	}
}

func TestSetDefaultString(t *testing.T) {
	tests := []struct {
		value        any
		defaultValue string
		expected     string
	}{
		{nil, "default", "default"},
		{"", "default", "default"},
		{"non-empty", "default", "non-empty"},
		{12345, "default", "12345"},
		{true, "default", "true"},
	}

	for _, tt := range tests {
		t.Run("SetDefaultString", func(t *testing.T) {
			result := SetDefaultString(tt.value, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %s, got %s for input %v", tt.expected, result, tt.value)
			}
		})
	}
}
