package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/router"
	pUtil "allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/allaboutapps/integresql-client-go"
	"github.com/allaboutapps/integresql-client-go/pkg/util"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var (
	client *integresql.Client
	hash   string

	// tracks template testDatabase initialization
	doOnce sync.Once

	migDir  = filepath.Join(pUtil.GetProjectRootDir(), "/migrations")
	fixFile = filepath.Join(pUtil.GetProjectRootDir(), "/internal/test/fixtures.go")
)

// Use this utility func to test with an isolated test database
func WithTestDatabase(t *testing.T, closure func(db *sql.DB)) {

	t.Helper()

	// new context derived from background
	ctx := context.Background()

	doOnce.Do(func() {

		t.Helper()
		initializeTestDatabaseTemplate(ctx, t)
	})

	testDatabase, err := client.GetTestDatabase(ctx, hash)

	if err != nil {
		t.Fatalf("Failed to obtain test database: %v", err)
	}

	connectionString := testDatabase.Config.ConnectionString()

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		t.Fatalf("Failed to setup test database for connectionString %q: %v", connectionString, err)
	}

	// this database object is managed and should close automatically after running the test
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping test database for connectionString %q: %v", connectionString, err)
	}

	t.Logf("WithTestDatabase: %q", testDatabase.Config.Database)

	closure(db)
}

// Use this utility func to test with an full blown server
func WithTestServer(t *testing.T, closure func(s *api.Server)) {

	t.Helper()

	WithTestDatabase(t, func(db *sql.DB) {

		t.Helper()

		defaultConfig := api.DefaultServiceConfigFromEnv()

		// https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
		// You may use port 0 to indicate you're not specifying an exact port but you want a free, available port selected by the system
		defaultConfig.Echo.ListenAddress = ":0"

		s := api.NewServer(defaultConfig)

		// attach the already initalized db
		s.DB = db

		if err := s.InitMailer(true); err != nil {
			t.Fatalf("failed to initialize mailer: %v", err)
		}

		router.Init(s)

		// no need to actually start echo!
		// see https://github.com/labstack/echo/issues/659

		closure(s)
	})
}

type GenericPayload map[string]interface{}

func (g GenericPayload) Reader(t *testing.T) *bytes.Reader {
	t.Helper()

	b, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("failed to serialize payload: %v", err)
	}

	return bytes.NewReader(b)
}

func PerformRequest(t *testing.T, s *api.Server, method string, path string, body GenericPayload, headers http.Header) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	if body == nil {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, body.Reader(t))
	}

	if headers != nil {
		req.Header = headers
	}
	if len(req.Header.Get(echo.HeaderContentType)) == 0 {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	res := httptest.NewRecorder()

	s.Echo.ServeHTTP(res, req)

	return res
}

func ParseResponseBody(t *testing.T, res *httptest.ResponseRecorder, v interface{}) {
	t.Helper()

	if err := json.NewDecoder(res.Result().Body).Decode(&v); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
}

func ParseResponseAndValidate(t *testing.T, res *httptest.ResponseRecorder, v interface{}) {
	t.Helper()

	ParseResponseBody(t, res, &v)

	val, ok := v.(runtime.Validatable)
	if !ok {
		t.Fatalf("Cannot parse response and validate, v (type %T) does not implement interface `runtime.Validatable`", v)
	}

	if err := val.Validate(strfmt.Default); err != nil {
		t.Fatalf("Failed to validate response: %v", err)
	}
}

func HeadersWithAuth(t *testing.T, token string) http.Header {
	t.Helper()

	headers := http.Header{}
	headers.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

	return headers
}

// main private function to properly build up the template database
// ensure it is called once once per pkg scope.
func initializeTestDatabaseTemplate(ctx context.Context, t *testing.T) {

	t.Helper()

	initTestDatabaseHash(t)

	initIntegresClient(t)

	if err := client.SetupTemplateWithDBClient(ctx, hash, func(db *sql.DB) error {

		t.Helper()

		err := applyMigrations(t, db)

		if err != nil {
			return err
		}

		err = insertFixtures(ctx, t, db)

		return err
	}); err != nil {

		// This error is exceptionally fatal as it hinders ANY future other
		// test execution with this hash as the template was *never* properly
		// setuped successfully. All GetTestDatabase will timeout.
		// TODO: Allow us to fail way faster here by telling the server that
		// this very template database is broken and cannot be used.

		t.Fatalf("Failed to setup template database for hash %q: %v", hash, err)
	}
}

func initIntegresClient(t *testing.T) {

	t.Helper()

	c, err := integresql.DefaultClientFromEnv()
	if err != nil {
		t.Fatalf("Failed to create new integresql-client: %v", err)
	}

	client = c
}

func initTestDatabaseHash(t *testing.T) {

	t.Helper()

	h, err := util.GetTemplateHash(migDir, fixFile)
	if err != nil {
		t.Fatalf("Failed to get template hash: %#v", err)
	}

	hash = h
}

func applyMigrations(t *testing.T, db *sql.DB) error {

	t.Helper()

	migrations := &migrate.FileMigrationSource{Dir: migDir}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	t.Logf("Applied %d migrations for hash %q", n, hash)

	return nil
}

func insertFixtures(ctx context.Context, t *testing.T, db *sql.DB) error {

	t.Helper()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	inserts := Inserts()

	for _, fixture := range inserts {
		if err := fixture.Insert(ctx, db, boil.Infer()); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	t.Logf("Inserted %d fixtures for hash %q", len(inserts), hash)

	return nil
}
