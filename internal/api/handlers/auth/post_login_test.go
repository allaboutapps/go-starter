package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"allaboutapps.at/aw/go-mranftl-sample/internal/api"
	"allaboutapps.at/aw/go-mranftl-sample/internal/api/handlers/auth"
	"allaboutapps.at/aw/go-mranftl-sample/internal/test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSuccessAuth(t *testing.T) {

	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {
		// cyclic import mini-test
		t.Log(auth.SAMPLE_EXPORTED_PGK_CONST)

		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
			"password": "password",
		}

		req := httptest.NewRequest("POST", "/api/v1/auth/login", payload.Reader(t))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		s.Echo.ServeHTTP(res, req)

		assert.Equal(t, http.StatusAccepted, res.Result().StatusCode)
	})

}

func TestInvalidCredentials(t *testing.T) {

	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {

		fixtures := test.Fixtures()
		payload := test.GenericPayload{
			"username": fixtures.User1.Username,
			"password": "not my password",
		}

		req := httptest.NewRequest("POST", "/api/v1/auth/login", payload.Reader(t))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		s.Echo.ServeHTTP(res, req)

		if res.Result().StatusCode != 401 {
			t.Logf("Did receive unexpected status code: %v", res.Result().StatusCode)
		}

	})

}
