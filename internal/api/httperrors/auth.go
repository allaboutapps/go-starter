package httperrors

import (
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/types"
)

var (
	ErrForbiddenUserDeactivated  = NewHTTPError(http.StatusForbidden, types.PublicHTTPErrorTypeUSERDEACTIVATED, "User account is deactivated")
	ErrBadRequestInvalidPassword = NewHTTPErrorWithDetail(http.StatusBadRequest, types.PublicHTTPErrorTypeINVALIDPASSWORD, "The password provided was invalid", "Password was either too weak or did not match other criteria")
	ErrForbiddenNotLocalUser     = NewHTTPError(http.StatusForbidden, types.PublicHTTPErrorTypeNOTLOCALUSER, "User account is not valid for local authentication")
	ErrNotFoundTokenNotFound     = NewHTTPError(http.StatusNotFound, types.PublicHTTPErrorTypeTOKENNOTFOUND, "Provided token was not found")
	ErrConflictTokenExpired      = NewHTTPError(http.StatusConflict, types.PublicHTTPErrorTypeTOKENEXPIRED, "Provided token has expired and is no longer valid")
	ErrConflictUserAlreadyExists = NewHTTPError(http.StatusConflict, types.PublicHTTPErrorTypeUSERALREADYEXISTS, "User with given username already exists")
)
