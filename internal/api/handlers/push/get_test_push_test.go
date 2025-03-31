package push_test

import (
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestGetTestPush(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		fix := fixtures.Fixtures()

		res := test.PerformRequest(t, s, "GET", "/api/v1/push/test", nil, test.HeadersWithAuth(t, fix.User1AccessToken1.Token))
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}

func TestGetTestPushUnauthorized(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/api/v1/push/test", nil, nil)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
	})
}
