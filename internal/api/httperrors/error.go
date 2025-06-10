package httperrors

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
)

// Payload in accordance with RFC 7807 (Problem Details for HTTP APIs) with the exception of the type
// value not being represented by a URI. https://tools.ietf.org/html/rfc7807 @ 2020-04-27T15:44:37Z

type HTTPError struct {
	types.PublicHTTPError
	Internal       error                  `json:"-"`
	AdditionalData map[string]interface{} `json:"-"`
}

type HTTPValidationError struct {
	types.PublicHTTPValidationError
	Internal       error                  `json:"-"`
	AdditionalData map[string]interface{} `json:"-"`
}

func NewHTTPError(code int, errorType types.PublicHTTPErrorType, title string) *HTTPError {
	return &HTTPError{
		PublicHTTPError: types.PublicHTTPError{
			Code:  swag.Int64(int64(code)),
			Type:  errorType.Pointer(),
			Title: swag.String(title),
		},
	}
}

func NewHTTPErrorWithDetail(code int, errorType types.PublicHTTPErrorType, title string, detail string) *HTTPError {
	return &HTTPError{
		PublicHTTPError: types.PublicHTTPError{
			Code:   swag.Int64(int64(code)),
			Type:   errorType.Pointer(),
			Title:  swag.String(title),
			Detail: detail,
		},
	}
}

func NewFromEcho(e *echo.HTTPError) *HTTPError {
	return NewHTTPError(e.Code, types.PublicHTTPErrorTypeGeneric, http.StatusText(e.Code))
}

func (e *HTTPError) Error() string {
	var b strings.Builder

	fmt.Fprintf(&b, "HTTPError %d (%s): %s", *e.Code, *e.Type, *e.Title)

	if len(e.Detail) > 0 {
		fmt.Fprintf(&b, " - %s", e.Detail)
	}
	if e.Internal != nil {
		fmt.Fprintf(&b, ", %v", e.Internal)
	}
	if len(e.AdditionalData) > 0 {
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

func NewHTTPValidationError(code int, errorType types.PublicHTTPErrorType, title string, validationErrors []*types.HTTPValidationErrorDetail) *HTTPValidationError {
	return &HTTPValidationError{
		PublicHTTPValidationError: types.PublicHTTPValidationError{
			PublicHTTPError: types.PublicHTTPError{
				Code:  swag.Int64(int64(code)),
				Type:  errorType.Pointer(),
				Title: swag.String(title),
			},
			ValidationErrors: validationErrors,
		},
	}
}

func NewHTTPValidationErrorWithDetail(code int, errorType types.PublicHTTPErrorType, title string, validationErrors []*types.HTTPValidationErrorDetail, detail string) *HTTPValidationError {
	return &HTTPValidationError{
		PublicHTTPValidationError: types.PublicHTTPValidationError{
			PublicHTTPError: types.PublicHTTPError{
				Code:   swag.Int64(int64(code)),
				Type:   errorType.Pointer(),
				Title:  swag.String(title),
				Detail: detail,
			},
			ValidationErrors: validationErrors,
		},
	}
}

func (e *HTTPValidationError) Error() string {
	var b strings.Builder

	fmt.Fprintf(&b, "HTTPValidationError %d (%s): %s", *e.Code, *e.Type, *e.Title)

	if len(e.Detail) > 0 {
		fmt.Fprintf(&b, " - %s", e.Detail)
	}
	if e.Internal != nil {
		fmt.Fprintf(&b, ", %v", e.Internal)
	}
	if len(e.AdditionalData) > 0 {
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

	b.WriteString(" - Validation: ")
	for i, ve := range e.ValidationErrors {
		fmt.Fprintf(&b, "%s (in %s): %s", *ve.Key, *ve.In, *ve.Error)
		if i < len(e.ValidationErrors)-1 {
			b.WriteString(", ")
		}
	}

	return b.String()
}
