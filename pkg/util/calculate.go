package util

import (
	"math"
	"strconv"
	"strings"
)

// CalculatePercentage calculates the result amount based on the input percentage and rounds it up.
func CalculatePercentage(amount, percentage int) int {
	result := float64(amount) * float64(percentage) / 100
	return int(math.Ceil(result))
}

// ToFormatRupiah formats an integer into Indonesian Rupiah currency format.
func ToFormatRupiah(amount int) string {
	// Convert the integer to a string
	number := strconv.Itoa(amount)

	// Reverse the string to handle thousands separator
	reversed := reverseString(number)

	// Add a dot every 3 digits
	var result strings.Builder
	for i, digit := range reversed {
		if i > 0 && i%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(digit)
	}

	// Reverse back to original order and add "Rp." prefix
	return "Rp." + reverseString(result.String())
}

// Helper function to reverse a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
