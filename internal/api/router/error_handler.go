package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
)

var (
	DefaultHTTPErrorHandlerConfig = HTTPErrorHandlerConfig{
		HideInternalServerErrorDetails: true,
	}
)

type HTTPErrorHandlerConfig struct {
	HideInternalServerErrorDetails bool
}

func HTTPErrorHandler() echo.HTTPErrorHandler {
	return HTTPErrorHandlerWithConfig(DefaultHTTPErrorHandlerConfig)
}

func HTTPErrorHandlerWithConfig(config HTTPErrorHandlerConfig) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var code int64
		var he error

		var httpError *httperrors.HTTPError
		var httpValidationError *httperrors.HTTPValidationError
		var echoHTTPError *echo.HTTPError

		switch {
		case errors.As(err, &httpError):
			code = *httpError.Code
			he = httpError

			if code == http.StatusInternalServerError && config.HideInternalServerErrorDetails {
				if httpError.Internal == nil {
					//nolint:errorlint
					httpError.Internal = fmt.Errorf("Internal Error: %s", httpError)
				}

				httpError.Title = swag.String(http.StatusText(http.StatusInternalServerError))
			}
		case errors.As(err, &httpValidationError):
			code = *httpValidationError.Code
			he = httpValidationError

			if code == http.StatusInternalServerError && config.HideInternalServerErrorDetails {
				if httpValidationError.Internal == nil {
					//nolint:errorlint
					httpValidationError.Internal = fmt.Errorf("Internal Error: %s", httpValidationError)
				}

				httpValidationError.Title = swag.String(http.StatusText(http.StatusInternalServerError))
			}
		case errors.As(err, &echoHTTPError):
			code = int64(echoHTTPError.Code)

			if code == http.StatusInternalServerError && config.HideInternalServerErrorDetails {
				if echoHTTPError.Internal == nil {
					//nolint:errorlint
					echoHTTPError.Internal = fmt.Errorf("Internal Error: %s", echoHTTPError)
				}

				he = &httperrors.HTTPError{
					PublicHTTPError: types.PublicHTTPError{
						Code:  swag.Int64(int64(echoHTTPError.Code)),
						Title: swag.String(http.StatusText(http.StatusInternalServerError)),
						Type:  swag.String(httperrors.HTTPErrorTypeGeneric),
					},
					Internal: echoHTTPError.Internal,
				}
			} else {
				msg, ok := echoHTTPError.Message.(string)
				if !ok {
					if m, errr := json.Marshal(msg); err == nil {
						msg = string(m)
					} else {
						msg = fmt.Sprintf("failed to marshal HTTP error message: %v", errr)
					}
				}

				he = &httperrors.HTTPError{
					PublicHTTPError: types.PublicHTTPError{
						Code:  swag.Int64(int64(echoHTTPError.Code)),
						Title: &msg,
						Type:  swag.String(httperrors.HTTPErrorTypeGeneric),
					},
					Internal: echoHTTPError.Internal,
				}
			}
		default:
			code = http.StatusInternalServerError
			if config.HideInternalServerErrorDetails {
				he = &httperrors.HTTPError{
					PublicHTTPError: types.PublicHTTPError{
						Code:  swag.Int64(int64(http.StatusInternalServerError)),
						Title: swag.String(http.StatusText(http.StatusInternalServerError)),
						Type:  swag.String(httperrors.HTTPErrorTypeGeneric),
					},

					Internal: err,
				}
			} else {
				he = &httperrors.HTTPError{
					PublicHTTPError: types.PublicHTTPError{
						Code:  swag.Int64(int64(http.StatusInternalServerError)),
						Title: swag.String(err.Error()),
						Type:  swag.String(httperrors.HTTPErrorTypeGeneric),
					},
				}
			}
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(int(code))
			} else {
				err = c.JSON(int(code), he)
			}

			if err != nil {
				util.LogFromEchoContext(c).Warn().Err(err).AnErr("http_err", err).Msg("Failed to handle HTTP error")
			}
		}
	}
}
