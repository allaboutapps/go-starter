// nolint:revive
package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-openapi/swag"
)

var (
	StringSpaceReplacer = regexp.MustCompile(`\s+`)
)

// GenerateRandomBytes returns n random bytes securely generated using the system's default CSPRNG.
//
// An error will be returned if reading from the secure random number generator fails, at which point
// the returned result should be discarded and not used any further.
func GenerateRandomBytes(n int) ([]byte, error) {
	result := make([]byte, n)

	_, err := rand.Read(result)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return result, nil
}

// GenerateRandomBase64String returns a string with n random bytes securely generated using the system's
// default CSPRNG in base64 encoding. The resulting string might not be of length n as the encoding for
// the raw bytes generated may vary.
//
// An error will be returned if reading from the secure random number generator fails, at which point
// the returned result should be discarded and not used any further.
func GenerateRandomBase64String(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// GenerateRandomHexString returns a string with n random bytes securely generated using the system's
// default CSPRNG in hexadecimal encoding. The resulting string might not be of length n as the encoding
// for the raw bytes generated may vary.
//
// An error will be returned if reading from the secure random number generator fails, at which point
// the returned result should be discarded and not used any further.
func GenerateRandomHexString(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

type CharRange int

const (
	CharRangeNumeric CharRange = iota
	CharRangeAlphaLowerCase
	CharRangeAlphaUpperCase
)

// GenerateRandomString returns a string with n random bytes securely generated using the system's
// default CSPRNG. The characters within the generated string will either be part of one ore more supplied
// range of characters, or based on characters in the extra string supplied.
//
// An error will be returned if reading from the secure random number generator fails, at which point
// the returned result should be discarded and not used any further.
func GenerateRandomString(n int, ranges []CharRange, extra string) (string, error) {
	var str strings.Builder

	if len(ranges) == 0 && len(extra) == 0 {
		return "", errors.New("random string can only be created if set of characters or extra string characters supplied")
	}

	validateFn := func(elem byte) bool {
		// IndexByte(string, byte) is basically Contains(string, string) without casting
		if strings.IndexByte(extra, elem) >= 0 {
			return true
		}

		for _, r := range ranges {
			switch r {
			case CharRangeNumeric:
				if elem >= '0' && elem <= '9' {
					return true
				}
			case CharRangeAlphaLowerCase:
				if elem >= 'a' && elem <= 'z' {
					return true
				}
			case CharRangeAlphaUpperCase:
				if elem >= 'A' && elem <= 'Z' {
					return true
				}
			}
		}

		return false
	}

	for str.Len() < n {
		buf, err := GenerateRandomBytes(n)
		if err != nil {
			return "", err
		}

		for _, b := range buf {
			if validateFn(b) {
				str.WriteByte(b)
			}
			if str.Len() >= n {
				break
			}
		}
	}

	return str.String(), nil
}

// Lowercases a string and trims whitespace from the beginning and end of the string
func ToUsernameFormat(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// NonEmptyOrNil returns a pointer to passed string if it is not empty. Passing empty strings returns nil instead.
func NonEmptyOrNil(s string) *string {
	if len(s) > 0 {
		return swag.String(s)
	}

	return nil
}

// EmptyIfNil returns an empty string if the passed pointer is nil. Passing a pointer to a string will return the value of the string.
func EmptyIfNil(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// ContainsAll returns true if a string (str) contains all substrings (sub).
func ContainsAll(str string, subs ...string) bool {
	subLen := len(subs)
	contains := make([]bool, subLen)
	indices := make([]int, subLen)
	substrings := make([][]rune, subLen)
	for i, substring := range subs {
		substrings[i] = []rune(substring)
	}

	for _, marked := range str {
		for i, sub := range substrings {
			if len(sub) == 0 {
				contains[i] = true
			}
			if !contains[i] && marked == sub[indices[i]] {
				indices[i]++
				if indices[i] >= len(sub) {
					contains[i] = true
				}
			}
		}
	}

	for _, c := range contains {
		if !c {
			return false
		}
	}

	return true
}
