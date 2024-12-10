package util

import "testing"

func TestEnvironmentTransform(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"local", "local"},
		{"dev", "development"},
		{"development", "development"},
		{"develop", "development"},
		{"stg", "staging"},
		{"staging", "staging"},
		{"production", "production"},
		{"prod", "production"},
		{"prd", "production"},
		{"unknown", ""},
		{"  dev  ", "development"},
		{"  staging  ", "staging"},
		{"PROD", "production"},
	}

	for _, tt := range tests {
		t.Run("EnvironmentTransform", func(t *testing.T) {
			result := EnvironmentTransform(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s for input %v", tt.expected, result, tt.input)
			}
		})
	}
}
