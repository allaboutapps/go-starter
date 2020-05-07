package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"allaboutapps.dev/aw/go-starter/internal/util/ref"
	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	var code int
	var he error

	switch e := err.(type) {
	case *types.HTTPError:
		code = *e.Code
		he = e
	case *types.HTTPValidationError:
		code = *e.Code
		he = e
	case *echo.HTTPError:
		msg, ok := e.Message.(string)
		if !ok {
			if m, errr := json.Marshal(msg); err == nil {
				msg = string(m)
			} else {
				msg = fmt.Sprintf("failed to marshal HTTP error message: %v", errr)
			}
		}

		code = e.Code
		he = &types.HTTPError{
			Code:     &e.Code,
			Title:    msg,
			Type:     types.HTTPErrorTypeGeneric,
			Internal: e.Internal,
		}
	default:
		code = http.StatusInternalServerError
		he = &types.HTTPError{
			Code:  ref.Int(http.StatusInternalServerError),
			Title: err.Error(),
			Type:  types.HTTPErrorTypeGeneric,
		}
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, he)
		}

		if err != nil {
			util.LogFromEchoContext(c).Warn().Err(err).AnErr("http_err", err).Msg("Failed to handle HTTP error")
		}
	}
}
