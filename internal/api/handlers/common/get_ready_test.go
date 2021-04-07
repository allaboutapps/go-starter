package common_test

import (
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetReadyReadiness(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/-/ready", nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
		require.Equal(t, res.Body.String(), "Ready.")
	})
}

func TestGetReadyReadinessBroken(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {

		// forcefully remove an initialized component to check if ready state works
		s.Mailer = nil

		res := test.PerformRequest(t, s, "GET", "/-/ready", nil, nil)
		require.Equal(t, 521, res.Result().StatusCode)
		require.Equal(t, "Not ready.", res.Body.String())
	})
}

func TestGetReadyDBBrokenNotReady(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {

		// forcefully remove pg
		err := s.DB.Close()
		require.NoError(t, err)

		res := test.PerformRequest(t, s, "GET", "/-/ready", nil, nil)
		require.Equal(t, 521, res.Result().StatusCode)
		require.Equal(t, "Not ready.", res.Body.String())
	})
}
