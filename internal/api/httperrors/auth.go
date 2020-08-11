package httperrors

import (
	"net/http"
)

var (
	ErrBadRequestInvalidPassword = NewHTTPErrorWithDetail(http.StatusBadRequest, "INVALID_PASSWORD", "The password provided was invalid", "Password was either too weak or did not match other criteria")
	ErrForbiddenNotLocalUser     = NewHTTPError(http.StatusForbidden, "NOT_LOCAL_USER", "User account is not valid for local authentication")
	ErrNotFoundTokenNotFound     = NewHTTPError(http.StatusNotFound, "TOKEN_NOT_FOUND", "Provided token was not found")
	ErrConflictTokenExpired      = NewHTTPError(http.StatusConflict, "TOKEN_EXPIRED", "Provided token has expired and is no longer valid")
	ErrConflictUserAlreadyExists = NewHTTPError(http.StatusConflict, "USER_ALREADY_EXISTS", "User with given username already exists")
)
