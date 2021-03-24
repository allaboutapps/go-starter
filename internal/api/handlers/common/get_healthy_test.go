package common_test

import (
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHealthySuccess(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {

		// explicitly set touchfile that no other test has (so we can explicitly remove it beforehand.)
		s.Config.Management.ProbeWriteableTouchfile = ".healthy-test"

		for _, writeablePath := range s.Config.Management.ProbeWriteablePathsAbs {
			os.Remove(path.Join(writeablePath, s.Config.Management.ProbeWriteableTouchfile))

			// also remove after test completion.
			defer os.Remove(path.Join(writeablePath, s.Config.Management.ProbeWriteableTouchfile))
		}

		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		// fmt.Println(res.Body.String())
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
		require.Contains(t, res.Body.String(), "seq_health=1")

		firstTouchTime := make([]time.Time, len(s.Config.Management.ProbeWriteablePathsAbs))

		// expect a new touchfiles were written
		for _, writeablePath := range s.Config.Management.ProbeWriteablePathsAbs {
			filePath := path.Join(writeablePath, s.Config.Management.ProbeWriteableTouchfile)
			stat, err := os.Stat(filePath)
			require.NoErrorf(t, err, "Expected to have %v", filePath)

			firstTouchTime = append(firstTouchTime, stat.ModTime())
		}

		res = test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, http.StatusOK, res.Result().StatusCode)
		require.Contains(t, res.Body.String(), "seq_health=2")

		// expect touchfiles modTime was updated
		for i, writeablePath := range s.Config.Management.ProbeWriteablePathsAbs {
			filePath := path.Join(writeablePath, s.Config.Management.ProbeWriteableTouchfile)
			stat, err := os.Stat(filePath)
			require.NoErrorf(t, err, "Expected to have %v", filePath)

			assert.NotEqual(t, firstTouchTime[i], stat.ModTime())
		}

		// fmt.Println(res.Body.String())
	})
}

func TestGetHealthyNoAuth(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/-/healthy", nil, nil)
		require.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})
}

func TestGetHealthyWrongAuth(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {
		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret=i-have-no-idea-about-the-pass", nil, nil)
		require.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
	})
}

func TestGetHealthyDBPingError(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {

		// forcefully close the DB
		s.DB.Close()

		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, 521, res.Result().StatusCode)
	})
}

func TestGetHealthyDBSeqError(t *testing.T) {
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
	test.WithTestServer(t, func(s *api.Server) {

		s.Config.Management.ProbeWriteablePathsAbs = []string{"/this/path/does/not/exist"}

		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, 521, res.Result().StatusCode)
	})
}

func TestGetHealthyNotReady(t *testing.T) {
	test.WithTestServer(t, func(s *api.Server) {

		// forcefully remove an initialized component to check if ready state works
		s.Mailer = nil

		res := test.PerformRequest(t, s, "GET", "/-/healthy?mgmt-secret="+s.Config.Management.Secret, nil, nil)
		require.Equal(t, 521, res.Result().StatusCode)
	})
}
