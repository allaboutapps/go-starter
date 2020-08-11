package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

type GenericPayload map[string]interface{}

func (g GenericPayload) Reader(t *testing.T) *bytes.Reader {
	t.Helper()

	b, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("failed to serialize payload: %v", err)
	}

	return bytes.NewReader(b)
}

func PerformRequestWithParams(t *testing.T, s *api.Server, method string, path string, body GenericPayload, headers http.Header, queryParams map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	if body == nil {
		return PerformRequestWithRawBody(t, s, method, path, nil, headers, queryParams)
	}

	return PerformRequestWithRawBody(t, s, method, path, body.Reader(t), headers, queryParams)
}

func PerformRequestWithRawBody(t *testing.T, s *api.Server, method string, path string, body io.Reader, headers http.Header, queryParams map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, body)

	if headers != nil {
		req.Header = headers
	}
	if body != nil && len(req.Header.Get(echo.HeaderContentType)) == 0 {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	res := httptest.NewRecorder()

	s.Echo.ServeHTTP(res, req)

	return res
}

func PerformRequest(t *testing.T, s *api.Server, method string, path string, body GenericPayload, headers http.Header) *httptest.ResponseRecorder {
	t.Helper()

	return PerformRequestWithParams(t, s, method, path, body, headers, nil)
}

func ParseResponseBody(t *testing.T, res *httptest.ResponseRecorder, v interface{}) {
	t.Helper()

	if err := json.NewDecoder(res.Result().Body).Decode(&v); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
}

func ParseResponseAndValidate(t *testing.T, res *httptest.ResponseRecorder, v runtime.Validatable) {
	t.Helper()

	ParseResponseBody(t, res, &v)

	if err := v.Validate(strfmt.Default); err != nil {
		t.Fatalf("Failed to validate response: %v", err)
	}
}

func HeadersWithAuth(t *testing.T, token string) http.Header {
	t.Helper()

	headers := http.Header{}
	headers.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

	return headers
}
