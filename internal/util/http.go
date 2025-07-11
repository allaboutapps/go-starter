package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"allaboutapps.dev/aw/go-starter/internal/api/httperrors"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/gabriel-vasile/mimetype"
	oerrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
)

const (
	HTTPHeaderCacheControl = "Cache-Control"
)

// BindAndValidateBody binds the request, parsing **only** its body (depending on the `Content-Type` request header) and performs validation
// as enforced by the Swagger schema associated with the provided type.
//
// Note: In contrast to BindAndValidate, this method does not restore the body after binding (it's considered consumed).
// Thus use BindAndValidateBody only once per request!
//
// Returns an error that can directly be returned from an echo handler and sent to the client should binding or validating of any model fail.
func BindAndValidateBody(c echo.Context, v runtime.Validatable) error {
	binder, ok := c.Echo().Binder.(*echo.DefaultBinder)
	if !ok {
		return fmt.Errorf("failed to get binder as *echo.DefaultBinder, got %T", c.Echo().Binder)
	}

	if err := binder.BindBody(c, v); err != nil {
		return fmt.Errorf("failed to bind body: %w", err)
	}

	return validatePayload(c, v)
}

// BindAndValidatePathAndQueryParams binds the request, parsing **only** its path **and** query params and performs validation
// as enforced by the Swagger schema associated with the provided type.
//
// Returns an error that can directly be returned from an echo handler and sent to the client should binding or validating of any model fail.
func BindAndValidatePathAndQueryParams(c echo.Context, v runtime.Validatable) error {
	binder, ok := c.Echo().Binder.(*echo.DefaultBinder)
	if !ok {
		return fmt.Errorf("failed to get binder as *echo.DefaultBinder, got %T", c.Echo().Binder)
	}

	if err := binder.BindPathParams(c, v); err != nil {
		return fmt.Errorf("failed to bind path params: %w", err)
	}

	if err := binder.BindQueryParams(c, v); err != nil {
		return fmt.Errorf("failed to bind query params: %w", err)
	}

	return validatePayload(c, v)
}

// BindAndValidatePathParams binds the request, parsing **only** its path params and performs validation
// as enforced by the Swagger schema associated with the provided type.
//
// Returns an error that can directly be returned from an echo handler and sent to the client should binding or validating of any model fail.
func BindAndValidatePathParams(c echo.Context, v runtime.Validatable) error {
	binder, ok := c.Echo().Binder.(*echo.DefaultBinder)
	if !ok {
		return fmt.Errorf("failed to get binder as *echo.DefaultBinder, got %T", c.Echo().Binder)
	}

	if err := binder.BindPathParams(c, v); err != nil {
		return fmt.Errorf("failed to bind path params: %w", err)
	}

	return validatePayload(c, v)
}

// BindAndValidateQueryParams binds the request, parsing **only** its query params and performs validation
// as enforced by the Swagger schema associated with the provided type.
//
// Returns an error that can directly be returned from an echo handler and sent to the client should binding or validating of any model fail.
func BindAndValidateQueryParams(c echo.Context, v runtime.Validatable) error {
	binder, ok := c.Echo().Binder.(*echo.DefaultBinder)
	if !ok {
		return fmt.Errorf("failed to get binder as *echo.DefaultBinder, got %T", c.Echo().Binder)
	}

	if err := binder.BindQueryParams(c, v); err != nil {
		return fmt.Errorf("failed to bind query params: %w", err)
	}

	return validatePayload(c, v)
}

// BindAndValidate binds the request, parsing path+query+body and validating these structs.
//
// Deprecated: Use our dedicated BindAndValidate* mappers instead:
//
//	BindAndValidateBody(c echo.Context, v runtime.Validatable) error // preferred
//	BindAndValidatePathAndQueryParams(c echo.Context, v runtime.Validatable) error  // preferred
//	BindAndValidatePathParams(c echo.Context, v runtime.Validatable) error // rare usecases
//	BindAndValidateQueryParams(c echo.Context, v runtime.Validatable) error // rare usecases
//
// BindAndValidate works like Echo <v4.2.0. It was preferred to .Bind() everything (query, params, body) to a single struct
// in one pass. Thus we included additional handling to allow multiple body rebindings (though copying while restoring),
// as goswagger generated structs per endpoint are typically **separated** into one params struct (path and query) and one
// body struct. Echo >=v4.2.0 DefaultBinder now supports binding query, path params and body to their **own** structs natively.
// Thus, you areencouraged to use our new dedicated BindAndValidate* mappers, which are relevant for the structs goswagger
// autogenerates for you.
//
// Original: Parses body (depending on the `Content-Type` request header) and performs payload validation as enforced by
// the Swagger schema associated with the provided type. In addition to binding the body, BindAndValidate also assigns query
// and URL parameters (if any) to a struct and perform validations on those.
//
// Providing more than one struct allows for binding payload and parameters simultaneously since echo and goswagger expect data
// to be structured differently. If you do not require parsing of both body and params, additional structs can be omitted.
//
// Returns an error that can directly be returned from an echo handler and sent to the client should binding or validating of any model fail.
func BindAndValidate(c echo.Context, v runtime.Validatable, validatables ...runtime.Validatable) error {
	// TODO error handling for all occurrences of Bind() due to JSON unmarshal type mismatches
	if len(validatables) == 0 {
		if err := defaultEchoBindAll(c, v); err != nil {
			return err
		}

		return validatePayload(c, v)
	}

	var reqBody []byte
	var err error
	if c.Request().Body != nil {
		reqBody, err = io.ReadAll(c.Request().Body)
		if err != nil {
			return fmt.Errorf("failed to read request body: %w", err)
		}
	}

	if err = restoreBindAndValidate(c, reqBody, v); err != nil {
		return fmt.Errorf("failed to restore bind and validate: %w", err)
	}

	for _, vv := range validatables {
		if err = restoreBindAndValidate(c, reqBody, vv); err != nil {
			return fmt.Errorf("failed to restore bind and validate: %w", err)
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

	if err := c.JSON(code, v); err != nil {
		return fmt.Errorf("failed to return JSON: %w", err)
	}

	return nil
}

func ParseFileUpload(c echo.Context, formNameFile string, allowedMIMETypes []string) (*multipart.FileHeader, multipart.File, *mimetype.MIME, error) {
	log := LogFromEchoContext(c)

	fileHeader, err := c.FormFile(formNameFile)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get form file")
		return nil, nil, nil, fmt.Errorf("failed to get form file: %w", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Debug().Err(err).Str("filename", fileHeader.Filename).Int64("fileSize", fileHeader.Size).Msg("Failed to open uploaded file")
		return nil, nil, nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}

	if fileHeader.Size < 1 {
		log.Debug().Err(err).Str("filename", fileHeader.Filename).Int64("fileSize", fileHeader.Size).Msg("File size can't be 0")
		return nil, nil, nil, httperrors.ErrBadRequestZeroFileSize
	}

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		log.Debug().Err(err).Str("filename", fileHeader.Filename).Int64("fileSize", fileHeader.Size).Msg("Failed to detect MIME type of uploaded file")
		file.Close()
		return nil, nil, nil, fmt.Errorf("failed to detect MIME type of uploaded file: %w", err)
	}

	// ! Important: we *MUST* reset the reader back to 0, since `minetype.DetectReader` reads the beginning of the
	// ! file in order to detect it's MIME type. Continuing to use the reader without resetting it results in a
	// ! corrupted file unable to be processed or opened otherwise.
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		log.Debug().Err(err).Str("filename", fileHeader.Filename).Int64("fileSize", fileHeader.Size).Msg("Failed to reset reader of uploaded file to start")
		file.Close()
		return nil, nil, nil, fmt.Errorf("failed to reset reader of uploaded file to start: %w", err)
	}

	allowed := false
	for _, allowedType := range allowedMIMETypes {
		if mime.Is(allowedType) {
			log.Debug().
				Str("mimeType", mime.String()).
				Str("mimeTypeFileExtension", mime.Extension()).
				Str("filename", fileHeader.Filename).
				Int64("fileSize", fileHeader.Size).
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
			Str("filename", fileHeader.Filename).
			Int64("fileSize", fileHeader.Size).
			Msg("MIME type of uploaded file is not allowed, rejecting")
		file.Close()

		return nil, nil, nil, echo.ErrUnsupportedMediaType
	}

	return fileHeader, file, mime, nil
}

func StreamFile(c echo.Context, code int, mediaType string, fileName string, r io.ReadCloser) error {
	formattedMediaType := mime.FormatMediaType("attachment",
		map[string]string{
			"filename": fileName,
		},
	)

	c.Response().Header().Set(echo.HeaderContentDisposition, formattedMediaType)
	SetOrAppendHeader(c.Response().Header(), echo.HeaderAccessControlExposeHeaders, echo.HeaderContentDisposition)

	defer r.Close()

	if err := c.Stream(code, mediaType, r); err != nil {
		return fmt.Errorf("failed to stream file: %w", err)
	}

	return nil
}

func SetOrAppendHeader(header http.Header, key string, values ...string) {
	headerSet := header.Get(key)
	headerValue := strings.Join(values, ", ")
	if headerSet == "" {
		header.Add(key, headerValue)
	} else {
		header.Set(key, strings.Join([]string{headerSet, headerValue}, ", "))
	}
}

func restoreBindAndValidate(c echo.Context, reqBody []byte, v runtime.Validatable) error {
	if reqBody != nil {
		c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	if err := defaultEchoBindAll(c, v); err != nil {
		return err
	}

	return validatePayload(c, v)
}

func validatePayload(c echo.Context, v runtime.Validatable) error {
	if err := v.Validate(strfmt.Default); err != nil {
		var compositeError *oerrors.CompositeError
		if errors.As(err, &compositeError) {
			LogFromEchoContext(c).Debug().Errs("validation_errors", compositeError.Errors).Msg("Payload did match schema, returning HTTP validation error")

			valErrs := formatValidationErrors(c.Request().Context(), compositeError)

			return httperrors.NewHTTPValidationError(http.StatusBadRequest, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)
		}

		var validationError *oerrors.Validation
		if errors.As(err, &validationError) {
			LogFromEchoContext(c).Debug().AnErr("validation_error", validationError).Msg("Payload did match schema, returning HTTP validation error")

			valErrs := []*types.HTTPValidationErrorDetail{
				{
					Key:   &validationError.Name,
					In:    &validationError.In,
					Error: swag.String(validationError.Error()),
				},
			}

			return httperrors.NewHTTPValidationError(http.StatusBadRequest, types.PublicHTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)
		}

		LogFromEchoContext(c).Error().Err(err).Msg("Failed to validate payload, returning generic HTTP error")

		return fmt.Errorf("failed to validate payload: %w", err)
	}

	return nil
}

// Bind it all
// Restores echo query binding pre 4.2.0 handling
// Newer echo versions no longer automatically bind query params to tagged :query struct-fields unless its a GET or DELETE request
// Workaround, depends on the internal echo.DefaultBinder methods.
//
// TODO: Eventually move to a customly implemented Binder.
// Hopefully BindPathParams, BindQueryParams and BindBody stay provided in the future.
//
// This upstream security fix does not directly affect us, as our goswagger generated params/query structs
// and body structs are separated from each other and cannot collide/overwrite props.
// https://github.com/labstack/echo/commit/4d626c210d3946814a30d545adf9b8f2296686a7#diff-aade326d3512b5a2ada6faa791ddec468f2a0adedb352339c9e314e74c8949d2
func defaultEchoBindAll(c echo.Context, v runtime.Validatable) error {
	binder, ok := c.Echo().Binder.(*echo.DefaultBinder)
	if !ok {
		return fmt.Errorf("failed to get binder as *echo.DefaultBinder, got %T", c.Echo().Binder)
	}

	if err := binder.BindPathParams(c, v); err != nil {
		return fmt.Errorf("failed to bind path params: %w", err)
	}
	if err := binder.BindQueryParams(c, v); err != nil {
		return fmt.Errorf("failed to bind query params: %w", err)
	}

	if err := binder.BindBody(c, v); err != nil {
		return fmt.Errorf("failed to bind body: %w", err)
	}

	return nil
}

func formatValidationErrors(ctx context.Context, compositeError *oerrors.CompositeError) []*types.HTTPValidationErrorDetail {
	valErrs := make([]*types.HTTPValidationErrorDetail, 0, len(compositeError.Errors))
	for _, err := range compositeError.Errors {
		var compositeError *oerrors.CompositeError
		if errors.As(err, &compositeError) {
			valErrs = append(valErrs, formatValidationErrors(ctx, compositeError)...)
			continue
		}

		var validationError *oerrors.Validation
		if errors.As(err, &validationError) {
			valErrs = append(valErrs, &types.HTTPValidationErrorDetail{
				Key:   &validationError.Name,
				In:    &validationError.In,
				Error: swag.String(validationError.Error()),
			})

			continue
		}

		LogFromContext(ctx).Warn().Err(err).Str("err_type", fmt.Sprintf("%T", err)).Msg("Received unknown error type while validating payload, skipping")
	}

	return valErrs
}
