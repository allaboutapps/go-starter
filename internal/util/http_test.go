package util_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/types/auth"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/strfmt/conv"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBindAndValidateSuccess(t *testing.T) {
	e := echo.New()
	//nolint:gosec
	testToken := "a546daf5-c845-46a7-8fa6-3d94ae7e1424"
	testResponse := &types.PostLoginResponse{
		AccessToken:  conv.UUID4(strfmt.UUID4("afbcbc30-4794-48bd-93f1-08373a031fe3")),
		RefreshToken: conv.UUID4(strfmt.UUID4("1dd1228c-fa9a-4755-b995-30e24dd6247d")),
		ExpiresIn:    swag.Int64(3600),
		TokenType:    swag.String("Bearer"),
	}

	e.POST("/", func(c echo.Context) error {
		testParam1 := auth.NewGetUserInfoRouteParams()
		testParam2 := auth.NewPostForgotPasswordRouteParams()
		var body types.PostRefreshPayload

		err := util.BindAndValidate(c, &body, &testParam1, &testParam2)
		assert.NoError(t, err)
		assert.NotEmpty(t, body)
		assert.Equal(t, strfmt.UUID4(testToken), *body.RefreshToken)

		return util.ValidateAndReturn(c, 200, testResponse)
	})
	testBody := test.GenericPayload{
		"refresh_token": testToken,
	}

	s := &api.Server{
		Echo: e,
	}

	res := test.PerformRequest(t, s, "POST", "/?test=true", testBody, nil)

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)

	var response types.PostLoginResponse
	test.ParseResponseAndValidate(t, res, &response)

	assert.Equal(t, *testResponse, response)
}

func TestBindAndValidateBadRequest(t *testing.T) {
	e := echo.New()
	testToken := "foo"

	e.POST("/", func(c echo.Context) error {
		var body types.PostRefreshPayload

		err := util.BindAndValidateBody(c, &body)
		assert.Error(t, err)

		return nil
	})
	testBody := test.GenericPayload{
		"refresh_token": testToken,
	}

	s := &api.Server{
		Echo: e,
	}

	_ = test.PerformRequest(t, s, "POST", "/?test=true", testBody, nil)
}

func TestParseFileUplaod(t *testing.T) {
	originalDocumentPath := filepath.Join(util.GetProjectRootDir(), "test", "testdata", "example.jpg")
	body, contentType := prepareFileUpload(t, originalDocumentPath)

	e := echo.New()
	e.POST("/", func(c echo.Context) error {

		fh, file, mime, err := util.ParseFileUpload(c, "file", []string{"image/jpeg"})
		require.NoError(t, err)
		assert.True(t, mime.Is("image/jpeg"))
		assert.NotEmpty(t, fh)
		assert.NotEmpty(t, file)

		return c.NoContent(204)
	})

	s := &api.Server{
		Echo: e,
	}

	headers := http.Header{}
	headers.Set(echo.HeaderContentType, contentType)

	res := test.PerformRequestWithRawBody(t, s, "POST", "/", body, headers, nil)

	require.Equal(t, http.StatusNoContent, res.Result().StatusCode)
}

func TestParseFileUplaodUnsupported(t *testing.T) {
	originalDocumentPath := filepath.Join(util.GetProjectRootDir(), "test", "testdata", "example.jpg")
	body, contentType := prepareFileUpload(t, originalDocumentPath)

	e := echo.New()
	e.POST("/", func(c echo.Context) error {

		fh, file, mime, err := util.ParseFileUpload(c, "file", []string{"image/png"})
		assert.Nil(t, fh)
		assert.Nil(t, file)
		assert.Nil(t, mime)
		if err != nil {
			return err
		}

		return c.NoContent(204)
	})

	s := &api.Server{
		Echo: e,
	}

	headers := http.Header{}
	headers.Set(echo.HeaderContentType, contentType)

	res := test.PerformRequestWithRawBody(t, s, "POST", "/", body, headers, nil)

	require.Equal(t, http.StatusUnsupportedMediaType, res.Result().StatusCode)
}

func prepareFileUpload(t *testing.T, filePath string) (*bytes.Buffer, string) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	src, err := os.Open(filePath)
	require.NoError(t, err)
	defer src.Close()

	dst, err := writer.CreateFormFile("file", src.Name())
	require.NoError(t, err)

	_, err = io.Copy(dst, src)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	return &body, writer.FormDataContentType()
}
