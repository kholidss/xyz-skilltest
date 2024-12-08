package util

import "github.com/spf13/cast"

func SetDefaultInt(value any, defaultValue int) int {
	number := cast.ToInt(value)

	if number == 0 {
		return defaultValue
	}

	return number
}

func SetDefaultString(value any, defaultValue string) string {
	str := cast.ToString(value)

	if str == "" {
		return defaultValue
	}

	return str
}
