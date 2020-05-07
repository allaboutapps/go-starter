package util

import (
	"fmt"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util/ref"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

var (
	ErrInternalValidationError = types.NewHTTPErrorWithDetail(http.StatusInternalServerError, "INTERNAL_VALIDATION_ERROR", "Payload validation failed", "Internal model used for payload validation does not implement `runtime.Validatable` interface")
)

// BindAndValidate binds the request, parsing either the body or form data (as defined by the `Content-Type` request header)
// and performs payload validation as enforced by the Swagger schema associated with the provided type.
// `v` must implement `github.com/go-openapi/runtime.Validatable` in order to perform validations, otherwise an internal server error is thrown.
// Returns an error that can directly be returned from an echo handler and sent to the client should binding or validating fail.
func BindAndValidate(c echo.Context, v interface{}) error {
	if err := c.Bind(&v); err != nil {
		return err
	}

	return validatePayload(c, v)
}

// ValidateAndReturn returns the provided data as a JSON response with the given HTTP status code after performing payload
// validation as enforced by the Swagger schema associated with the provided type.
// `v` must implement `github.com/go-openapi/runtime.Validatable` in order to perform validations, otherwise an internal server error is thrown.
// Returns an error that can directly be returned from an echo handler and sent to the client should sending or validating fail.
func ValidateAndReturn(c echo.Context, code int, v interface{}) error {
	if err := validatePayload(c, v); err != nil {
		return err
	}

	return c.JSON(code, v)
}

func validatePayload(c echo.Context, v interface{}) error {
	val, ok := v.(runtime.Validatable)
	if !ok {
		LogFromEchoContext(c).Error().Str("type", fmt.Sprintf("%T", v)).Msg("Type does not implement interface `runtime.Validatable`, cannot validate")
		return ErrInternalValidationError
	}

	if err := val.Validate(strfmt.Default); err != nil {
		compErr, ok := err.(*errors.CompositeError)
		if ok {
			LogFromEchoContext(c).Debug().Errs("validation_errors", compErr.Errors).Msg("Payload did match schema, returning HTTP validation error")

			valErrs := make([]*types.HTTPValidationErrorDetail, 0, len(compErr.Errors))
			for _, e := range compErr.Errors {
				if valErr, ok := e.(*errors.Validation); ok {
					valErrs = append(valErrs, &types.HTTPValidationErrorDetail{
						Key:   valErr.Name,
						In:    valErr.In,
						Error: valErr.Error(),
					})
				} else {
					LogFromEchoContext(c).Warn().Err(e).Str("err_type", fmt.Sprintf("%T", e)).Msg("Received unknown error type while validating payload, skipping")
				}
			}

			return &types.HTTPValidationError{
				HTTPError: types.HTTPError{
					Code:  ref.Int(http.StatusBadRequest),
					Type:  types.HTTPErrorTypeGeneric,
					Title: http.StatusText(http.StatusBadRequest),
				},
				ValidationErrors: valErrs,
			}
		}

		LogFromEchoContext(c).Error().Err(err).Msg("Failed to validate payload, returning generic HTTP error")
		return err
	}

	return nil
}
