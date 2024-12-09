package masker

import "strings"

func Censored(s, masked string) string {

	if len(masked) == 0 {
		masked = "*"
	}

	rs := []rune(s)
	replaceChar := []rune(masked)[0]
	const skipChar = '-'

	var start, end int

	if len(rs) <= 10 {
		start = 2
		end = len(rs) - 1
	}

	if len(rs) > 10 {
		start = 3
		end = len(rs) - 4
	}

	for i := start; i < end; i++ {
		if rs[i] != skipChar {
			rs[i] = replaceChar
		}
	}

	return string(rs)
}

func FullCensor(s, m string) string {
	return strings.Repeat(m, len(s))
}
