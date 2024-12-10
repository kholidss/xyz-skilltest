package util

import "testing"

func TestSanitizeDate(t *testing.T) {
	tests := []struct {
		input    string
		layout   string
		expected string
	}{
		{"2024-12-10", "2006-01-02", "2024-12-10"},
		{"10/12/2024", "2006-01-02", "2024-12-10"},
		{"12/10/2024", "2006-01-02", "2024-10-12"},
		{"10-12-2024", "2006-01-02", "2024-12-10"},
		{"12-10-2024", "2006-01-02", "2024-10-12"},
		{"2024.12.10", "2006-01-02", "2024-12-10"},
		{"Dec 10, 2024", "2006-01-02", "2024-12-10"},
		{"10 Dec 2024", "2006-01-02", "2024-12-10"},
		{"December 10, 2024", "2006-01-02", "2024-12-10"},
		{"10 December 2024", "2006-01-02", "2024-12-10"},
		{"invalid-date", "2006-01-02", "invalid-date"},
	}

	for _, tt := range tests {
		t.Run("SanitizeDate", func(t *testing.T) {
			result := SanitizeDate(tt.input, tt.layout)
			if result != tt.expected {
				t.Errorf("expected %s, got %s for input %v", tt.expected, result, tt.input)
			}
		})
	}
}
