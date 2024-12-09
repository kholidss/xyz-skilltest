package util

import (
	"time"
)

// SanitizeDate sanitize new date with layout format given
func SanitizeDate(input string, layout string) string {
	formats := []string{
		"2006-01-02",      // ISO format (default for Go's time.Parse)
		"02/01/2006",      // DD/MM/YYYY
		"01/02/2006",      // MM/DD/YYYY
		"02-01-2006",      // DD-MM-YYYY
		"01-02-2006",      // MM-DD-YYYY
		"2006.01.02",      // YYYY.MM.DD
		"Jan 2, 2006",     // Month day, year
		"2 Jan 2006",      // Day Month Year
		"January 2, 2006", // Full month name
		"2 January 2006",  // Day Full month name Year
	}

	for _, format := range formats {
		parsedDate, err := time.Parse(format, input)
		if err == nil {
			// Format the parsed date into "YYYY-MM-DD"
			return parsedDate.Format(layout)
		}
	}

	return input
}
