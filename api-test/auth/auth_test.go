package auth

import (
	"net/http/httptest"
	"strings"
	"testing"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/test"
	"github.com/labstack/echo/v4"
)

func TestSuccessAuth(t *testing.T) {

	t.Parallel()

	test.WithTestServer(func(s *api.Server) {

		userJSON := `{
			"username": "user1@example.com",
			"password": "password"
		}`

		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		s.Echo.ServeHTTP(res, req)

		if res.Result().StatusCode != 200 {
			t.Logf("Did receive unexpected status code: %v", res.Result().StatusCode)
		}

	})

	// WithTestDatabase(func(db *sql.DB) {
	// 	err := user1.Reload(context.Background(), db)

	// 	if err != nil {
	// 		t.Error("failed to reload")
	// 	}

	// 	// fmt.Println(user1)
	// })

}
