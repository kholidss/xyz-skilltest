package util

import "testing"

func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		amount     int
		percentage int
		expected   int
	}{
		{100, 10, 10},
		{250, 20, 50},
		{99, 50, 50},
		{500, 0, 0},
		{100, 99, 99},
		{50, 25, 13},
	}

	for _, tt := range tests {
		t.Run("CalculatePercentage", func(t *testing.T) {
			result := CalculatePercentage(tt.amount, tt.percentage)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestToFormatRupiah(t *testing.T) {
	tests := []struct {
		amount   int
		expected string
	}{
		{1000, "Rp.1.000"},
		{1000000, "Rp.1.000.000"},
		{500, "Rp.500"},
		{999999999, "Rp.999.999.999"},
		{123456789, "Rp.123.456.789"},
	}

	for _, tt := range tests {
		t.Run("ToFormatRupiah", func(t *testing.T) {
			result := ToFormatRupiah(tt.amount)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
