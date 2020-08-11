package db

import "strings"

func SearchStringToTSQuery(s string) string {
	return strings.ReplaceAll(strings.Trim(s, " "), " ", ":* & ") + ":*"
}
