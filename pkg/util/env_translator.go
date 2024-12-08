package util

import "strings"

var envArr = map[string]string{
	"local": "local",

	"dev":         "development",
	"development": "development",
	"develop":     "development",

	"stg":     "staging",
	"staging": "staging",

	"production": "production",
	"prod":       "production",
	"prd":        "production",
}

func EnvironmentTransform(s string) string {
	v, ok := envArr[strings.ToLower(strings.Trim(s, " "))]

	if !ok {
		return ""
	}

	return v
}
