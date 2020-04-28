package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"allaboutapps.at/aw/go-mranftl-sample/internal/util"
	"allaboutapps.at/aw/go-mranftl-sample/internal/types"
	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	he, ok := err.(*types.HTTPError)
	if !ok {
		if hee, ok := err.(*echo.HTTPError); ok {
			msg, ok := hee.Message.(string)
			if !ok {
				if m, errr := json.Marshal(msg); err == nil {
					msg = string(m)
				} else {
					msg = fmt.Sprintf("failed to marshal HTTP error message: %v", errr)
				}
			}

			he = &types.HTTPError{
				Code:     hee.Code,
				Title:    msg,
				Type:     types.HTTPErrorTypeGeneric,
				Internal: hee.Internal,
			}
		} else {
			he = &types.HTTPError{
				Code:  http.StatusInternalServerError,
				Title: err.Error(),
				Type:  types.HTTPErrorTypeGeneric,
			}
		}
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(he.Code, he)
		}

		if err != nil {
			util.LogFromEchoContext(c).Warn().Err(err).Msg("Failed to handle HTTP error")
		}
	}
}
