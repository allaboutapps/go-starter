package router

import (
	"encoding/json"
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

		switch e := err.(type) {
		case *httperrors.HTTPError:
			code = *e.Code
			he = e

			if code == http.StatusInternalServerError && config.HideInternalServerErrorDetails {
				if e.Internal == nil {
					e.Internal = fmt.Errorf("%s", e.Error())
				}

				e.Title = swag.String(http.StatusText(http.StatusInternalServerError))
			}
		case *httperrors.HTTPValidationError:
			code = *e.Code
			he = e

			if code == http.StatusInternalServerError && config.HideInternalServerErrorDetails {
				if e.Internal == nil {
					e.Internal = fmt.Errorf("%s", e.Error())
				}

				e.Title = swag.String(http.StatusText(http.StatusInternalServerError))
			}
		case *echo.HTTPError:
			code = int64(e.Code)

			if code == http.StatusInternalServerError && config.HideInternalServerErrorDetails {
				if e.Internal == nil {
					e.Internal = fmt.Errorf("%s", e.Error())
				}

				he = &httperrors.HTTPError{
					PublicHTTPError: types.PublicHTTPError{
						Code:  swag.Int64(int64(e.Code)),
						Title: swag.String(http.StatusText(http.StatusInternalServerError)),
						Type:  swag.String(httperrors.HTTPErrorTypeGeneric),
					},
					Internal: e.Internal,
				}
			} else {
				msg, ok := e.Message.(string)
				if !ok {
					if m, errr := json.Marshal(msg); err == nil {
						msg = string(m)
					} else {
						msg = fmt.Sprintf("failed to marshal HTTP error message: %v", errr)
					}
				}

				he = &httperrors.HTTPError{
					PublicHTTPError: types.PublicHTTPError{
						Code:  swag.Int64(int64(e.Code)),
						Title: &msg,
						Type:  swag.String(httperrors.HTTPErrorTypeGeneric),
					},
					Internal: e.Internal,
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

					Internal: e,
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
