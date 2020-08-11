package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
)

// BindAndValidate binds the request, parsing its body (depending on the `Content-Type` request header) and performs payload
// validation as enforced by the Swagger schema associated with the provided type. In addition to binding the body, BindAndValidate
// can also assign query and URL parameters to a struct and perform validations on those.
//
// Providing more than one struct allows for binding payload and parameters simultaneously since echo and goswagger expect data
// to be structured differently. If you do not require parsing of both body and params, additional structs can be omitted.
//
// Returns an error that can directly be returned from an echo handler and sent to the client should binding or validating of any model fail.
func BindAndValidate(c echo.Context, v runtime.Validatable, vs ...runtime.Validatable) error {
	// TODO error handling for all occurrences of Bind() due to JSON unmarshal type mismatches
	if len(vs) == 0 {
		if err := c.Bind(v); err != nil {
			return err
		}

		return validatePayload(c, v)
	}

	var reqBody []byte = nil
	var err error
	if c.Request().Body != nil {
		reqBody, err = ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
	}

	if err = restoreBindAndValidate(c, reqBody, v); err != nil {
		return err
	}

	for _, vv := range vs {
		if err = restoreBindAndValidate(c, reqBody, vv); err != nil {
			return err
		}
	}

	return nil
}

// ValidateAndReturn returns the provided data as a JSON response with the given HTTP status code after performing payload
// validation as enforced by the Swagger schema associated with the provided type.
// `v` must implement `github.com/go-openapi/runtime.Validatable` in order to perform validations, otherwise an internal server error is thrown.
// Returns an error that can directly be returned from an echo handler and sent to the client should sending or validating fail.
func ValidateAndReturn(c echo.Context, code int, v runtime.Validatable) error {
	if err := validatePayload(c, v); err != nil {
		return err
	}

	return c.JSON(code, v)
}

func ParseFileUpload(c echo.Context, formNameFile string, allowedMIMETypes []string) (*multipart.FileHeader, multipart.File, *mimetype.MIME, error) {
	log := LogFromEchoContext(c)

	fh, err := c.FormFile(formNameFile)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get form file")
		return nil, nil, nil, err
	}

	file, err := fh.Open()
	if err != nil {
		log.Debug().Err(err).Str("filename", fh.Filename).Int64("fileSize", fh.Size).Msg("Failed to open uploaded file")
		return nil, nil, nil, err
	}

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		log.Debug().Err(err).Str("filename", fh.Filename).Int64("fileSize", fh.Size).Msg("Failed to detect MIME type of uploaded file")
		file.Close()
		return nil, nil, nil, err
	}

	// ! Important: we *MUST* reset the reader back to 0, since `minetype.DetectReader` reads the beginning of the
	// ! file in order to detect it's MIME type. Continuing to use the reader without resetting it results in a
	// ! corrupted file unable to be processed or opened otherwise.
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		log.Debug().Err(err).Str("filename", fh.Filename).Int64("fileSize", fh.Size).Msg("Failed to reset reader of uploaded file to start")
		file.Close()
		return nil, nil, nil, err
	}

	allowed := false
	for _, allowedType := range allowedMIMETypes {
		if mime.Is(allowedType) {
			log.Debug().
				Str("mimeType", mime.String()).
				Str("mimeTypeFileExtension", mime.Extension()).
				Str("filename", fh.Filename).
				Int64("fileSize", fh.Size).
				Str("allowedMIMEType", allowedType).
				Msg("MIME type of uploaded file is allowed, processing")

			allowed = true
			break
		}
	}

	if !allowed {
		log.Debug().
			Str("mimeType", mime.String()).
			Str("mimeTypeFileExtension", mime.Extension()).
			Str("filename", fh.Filename).
			Int64("fileSize", fh.Size).
			Msg("MIME type of uploaded file is not allowed, rejecting")
		file.Close()
		return nil, nil, nil, echo.ErrUnsupportedMediaType
	}

	return fh, file, mime, nil
}

func restoreBindAndValidate(c echo.Context, reqBody []byte, v runtime.Validatable) error {
	if reqBody != nil {
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	}

	if err := c.Bind(v); err != nil {
		return err
	}

	return validatePayload(c, v)
}

func validatePayload(c echo.Context, v runtime.Validatable) error {
	if err := v.Validate(strfmt.Default); err != nil {
		switch ee := err.(type) {
		case *errors.CompositeError:
			LogFromEchoContext(c).Debug().Errs("validation_errors", ee.Errors).Msg("Payload did match schema, returning HTTP validation error")

			valErrs := formatValidationErrors(c.Request().Context(), ee)

			return httperrors.NewHTTPValidationError(http.StatusBadRequest, httperrors.HTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)
		case *errors.Validation:
			LogFromEchoContext(c).Debug().AnErr("validation_error", ee).Msg("Payload did match schema, returning HTTP validation error")

			valErrs := []*types.HTTPValidationErrorDetail{
				{
					Key:   &ee.Name,
					In:    &ee.In,
					Error: swag.String(ee.Error()),
				},
			}

			return httperrors.NewHTTPValidationError(http.StatusBadRequest, httperrors.HTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)
		default:
			LogFromEchoContext(c).Error().Err(err).Msg("Failed to validate payload, returning generic HTTP error")
			return err
		}
	}

	return nil
}

func formatValidationErrors(ctx context.Context, err *errors.CompositeError) []*types.HTTPValidationErrorDetail {
	valErrs := make([]*types.HTTPValidationErrorDetail, 0, len(err.Errors))
	for _, e := range err.Errors {
		switch ee := e.(type) {
		case *errors.Validation:
			valErrs = append(valErrs, &types.HTTPValidationErrorDetail{
				Key:   &ee.Name,
				In:    &ee.In,
				Error: swag.String(ee.Error()),
			})
		case *errors.CompositeError:
			valErrs = append(valErrs, formatValidationErrors(ctx, ee)...)
		default:
			LogFromContext(ctx).Warn().Err(e).Str("err_type", fmt.Sprintf("%T", e)).Msg("Received unknown error type while validating payload, skipping")
		}
	}

	return valErrs
}
