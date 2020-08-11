package httperrors

import (
	"net/http"
)

var (
	ErrConflictPushToken    = NewHTTPError(http.StatusConflict, "PUSH_TOKEN_ALREADY_EXISTS", "The given token already exists.")
	ErrNotFoundOldPushToken = NewHTTPError(http.StatusNotFound, "OLD_PUSH_TOKEN_NOT_FOUND", "The old push token does not exists. The new token was saved.")
)
