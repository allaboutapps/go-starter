package httperrors

import "net/http"

var (
	ErrBadRequestZeroFileSize = NewHTTPError(http.StatusBadRequest, "ZERO_FILE_SIZE", "File size of 0 is not supported.")
)
