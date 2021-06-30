package db

import (
	"regexp"
	"strings"
)

var (
	tsQueryWhiteSpaceRegex = regexp.MustCompile(`\s+`)
)

// SearchStringToTSQuery returns a TSQuery string from user input.
// The resulting query will match if every word matches a beginning of a word in the row.
// This function will trim all leading and trailing as well as consecutive whitespaces and remove all single quotes before
// transforming the input into TSQuery syntax.
// If no input was given (nil or empty string) or the value only contains invalid characters, an empty string will be returned.
func SearchStringToTSQuery(s *string) string {
	if s == nil || len(*s) == 0 {
		return ""
	}

	v := strings.TrimSpace(strings.ReplaceAll(*s, "'", ""))
	if len(v) == 0 {
		return ""
	}

	return "'" + tsQueryWhiteSpaceRegex.ReplaceAllString(v, "':* & '") + "':*"
}
