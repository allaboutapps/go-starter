package types

import (
	"fmt"
	"sort"
	"strings"
)

const (
	HTTPErrorTypeGeneric string = "generic"
)

// Payload in accordance to RFC 7807 (Problem Details for HTTP APIs) with the exception of the type
// value not being represented by a URI. https://tools.ietf.org/html/rfc7807 @ 2020-04-27T15:44:37Z

// swagger:model
type HTTPError struct {
	// HTTP status code returned for the error
	// min: 100
	// max: 599
	// example: 403
	Code int `json:"status"`
	// Type of error returned, should be used for client-side error handling
	// required: true
	// example: generic
	Type string `json:"type"`
	// Short, human-readable description of the error
	// required: true
	// example: Forbidden
	Title string `json:"title"`
	// More detailed, human-readable optional explanation of the error
	// example: User is lacking permission to access this resource
	Detail         string                 `json:"detail,omitempty"`
	Internal       error                  `json:"-"`
	AdditionalData map[string]interface{} `json:"-"`
}

func NewHTTPError(code int, errorType string, title string) *HTTPError {
	return &HTTPError{
		Code:  code,
		Type:  errorType,
		Title: title,
	}
}

func NewHTTPErrorWithDetail(code int, errorType string, title string, detail string) *HTTPError {
	return &HTTPError{
		Code:   code,
		Type:   errorType,
		Title:  title,
		Detail: detail,
	}
}

func (e *HTTPError) Error() string {
	var b strings.Builder

	fmt.Fprintf(&b, "HTTPError %d (%s): %s", e.Code, e.Type, e.Title)

	if len(e.Detail) > 0 {
		fmt.Fprintf(&b, " - %s", e.Detail)
	}
	if e.Internal != nil {
		fmt.Fprintf(&b, ", %v", e.Internal)
	}
	if e.AdditionalData != nil && len(e.AdditionalData) > 0 {
		keys := make([]string, 0, len(e.AdditionalData))
		for k := range e.AdditionalData {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		b.WriteString(". Additional: ")
		for i, k := range keys {
			fmt.Fprintf(&b, "%s=%v", k, e.AdditionalData[k])
			if i < len(keys)-1 {
				b.WriteString(", ")
			}
		}
	}

	return b.String()
}
