package db

import "strings"

// SearchStringToTSQuery returns a TSQuery string from user input.
// The resulting query will match if every word matches a beginning of a word in the row.
func SearchStringToTSQuery(s string) string {
	return strings.ReplaceAll(strings.Trim(s, " "), " ", ":* & ") + ":*"
}
