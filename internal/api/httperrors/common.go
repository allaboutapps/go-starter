package httperrors

import (
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/types"
)

var (
	ErrBadRequestZeroFileSize = NewHTTPError(http.StatusBadRequest, types.PublicHTTPErrorTypeZEROFILESIZE, "File size of 0 is not supported.")
)
