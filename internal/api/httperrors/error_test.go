package httperrors_test

import (
	"database/sql"
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/require"
)

func TestHTTPErrorSimple(t *testing.T) {
	err := httperrors.NewHTTPError(http.StatusNotFound, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusNotFound))
	require.Equal(t, "HTTPError 404 (generic): Not Found", err.Error())
}

func TestHTTPErrorDetail(t *testing.T) {
	err := httperrors.NewHTTPErrorWithDetail(http.StatusNotFound, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusNotFound), "ToS violation")
	require.Equal(t, "HTTPError 404 (generic): Not Found - ToS violation", err.Error())
}

func TestHTTPErrorInternalError(t *testing.T) {
	err := httperrors.NewHTTPError(http.StatusInternalServerError, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusInternalServerError))

	err.Internal = sql.ErrConnDone

	require.Equal(t, "HTTPError 500 (generic): Internal Server Error, sql: connection is already closed", err.Error())
}

func TestHTTPErrorAdditionalData(t *testing.T) {
	err := httperrors.NewHTTPError(http.StatusInternalServerError, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusInternalServerError))

	err.AdditionalData = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	require.Equal(t, "HTTPError 500 (generic): Internal Server Error. Additional: key1=value1, key2=value2", err.Error())
}

var valErrs = append(make([]*types.HTTPValidationErrorDetail, 0, 2), &types.HTTPValidationErrorDetail{
	Key:   swag.String("test1"),
	In:    swag.String("body.test1"),
	Error: swag.String("ValidationError"),
}, &types.HTTPValidationErrorDetail{
	Key:   swag.String("test2"),
	In:    swag.String("body.test2"),
	Error: swag.String("Validation Error"),
})

func TestHTTPValidationErrorSimple(t *testing.T) {
	err := httperrors.NewHTTPValidationError(http.StatusBadRequest, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)
	require.Equal(t, "HTTPValidationError 400 (generic): Bad Request - Validation: test1 (in body.test1): ValidationError, test2 (in body.test2): Validation Error", err.Error())
}

func TestHTTPValidationErrorDetail(t *testing.T) {
	err := httperrors.NewHTTPValidationErrorWithDetail(http.StatusBadRequest, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs, "Did API spec change?")
	require.Equal(t, "HTTPValidationError 400 (generic): Bad Request - Did API spec change? - Validation: test1 (in body.test1): ValidationError, test2 (in body.test2): Validation Error", err.Error())
}

func TestHTTPValidationErrorInternalError(t *testing.T) {
	err := httperrors.NewHTTPValidationError(http.StatusBadRequest, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)

	err.Internal = sql.ErrConnDone

	require.Equal(t, "HTTPValidationError 400 (generic): Bad Request, sql: connection is already closed - Validation: test1 (in body.test1): ValidationError, test2 (in body.test2): Validation Error", err.Error())
}

func TestHTTPValidationErrorAdditionalData(t *testing.T) {
	err := httperrors.NewHTTPValidationError(http.StatusBadRequest, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)

	err.AdditionalData = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	require.Equal(t, "HTTPValidationError 400 (generic): Bad Request. Additional: key1=value1, key2=value2 - Validation: test1 (in body.test1): ValidationError, test2 (in body.test2): Validation Error", err.Error())
}
