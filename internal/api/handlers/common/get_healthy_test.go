package common_test

import (
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/require"
)

func TestGetHealthySuccess(t *testing.T) {
	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
		require.Contains(t, res.Body.String(), "seq_health=1")

		res = test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
		require.Contains(t, res.Body.String(), "seq_health=2")

		// fmt.Println(res.Body.String())
	})
}

func TestGetHealthyNoAuth(t *testing.T) {
	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/-/healthy", nil, nil)
		require.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})
}

func TestGetHealthyWrongAuth(t *testing.T) {
	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret=i-have-no-idea-about-the-pass", nil, nil)
		require.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
	})
}

func TestGetHealthyDBPingError(t *testing.T) {
	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {

		// forcefully close the DB
		s.DB.Close()

		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, 521, res.Result().StatusCode)
	})
}

func TestGetHealthyDBSeqError(t *testing.T) {
	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {

		// forcefully remove the sequence
		if _, err := s.DB.Exec("DROP SEQUENCE seq_health;"); err != nil {
			t.Fatal(err, "was unable to drop sequence seq_health")
		}

		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)

		require.Equal(t, 521, res.Result().StatusCode)
	})
}

func TestGetHealthyMountError(t *testing.T) {
	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {

		s.Config.Paths.MntBaseDirAbs = "/this/path/does/not/exist"

		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, 521, res.Result().StatusCode)
	})
}

func TestGetHealthyNotReady(t *testing.T) {
	t.Parallel()

	test.WithTestServer(t, func(s *api.Server) {

		// forcefully remove an initialized component to check if ready state works
		s.Mailer = nil

		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, 521, res.Result().StatusCode)
	})
}
