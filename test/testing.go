package test

import (
	"context"
	"database/sql"
	"log"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/router"
	"allaboutapps.at/aw/go-mranftl-sample/test/pgconsumer"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/volatiletech/sqlboiler/boil"
)

var (
	client *pgconsumer.Client
	hash   string
	// ! TODO golang does not support relative paths in files properly
	// It's only possible to supply this by
	// Use ENV var to specify app-root
	migDir  = "/app/migrations"
	fixFile = "/app/test/fixtures.go"
)

func initIntegres() {

	c, err := pgconsumer.DefaultClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create new pgconsumer client: %v", err)
	}

	client = c
}

func initHash() {

	h, err := pgconsumer.GetTemplateHash(migDir, fixFile)
	if err != nil {
		log.Fatalf("Failed to get template hash: %#v", err)
	}

	hash = h
}

func initTemplate(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{Dir: migDir}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Printf("Applied %d migrations for hash %q", n, hash)

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, fixture := range fixtures {
		if err := fixture.Insert(context.Background(), db, boil.Infer()); err != nil {
			if errr := tx.Rollback(); errr != nil {
				return errr
			}

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Inserted %d fixtures for hash %q", len(fixtures), hash)

	return nil
}

func InitializeDatabaseTemplate() {

	initHash()

	initIntegres()

	if err := client.SetupTemplateWithDBClient(context.Background(), hash, initTemplate); err != nil {
		log.Fatalf("Failed to setup template database for hash %q: %v", hash, err)
	}
}

// Use this utility func to test with an isolated test database
func WithTestDatabase(closure func(db *sql.DB)) {
	testDatabase, err := client.GetTestDatabase(context.Background(), hash)

	if err != nil {
		log.Fatalf("Failed to obtain TestDatabase: %v", err)
	}

	connectionString := testDatabase.Config.ConnectionString()

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatalf("Failed to setup testdatabase for connectionString %q: %v", connectionString, err)
	}

	// this database object is managed and should close automatically after running the test
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("Failed to ping testdatabase for connectionString %q: %v", connectionString, err)
	}

	closure(db)
}

// Use this utility func to test with an full blown server
func WithTestServer(closure func(s *api.Server)) {
	WithTestDatabase(func(db *sql.DB) {

		defaultConfig := api.DefaultServiceConfigFromEnv()

		// https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
		// You may use port 0 to indicate you're not specifying an exact port but you want a free, available port selected by the system
		defaultConfig.Echo.ListenAddress = ":0"

		s := api.NewServer(defaultConfig)

		// attach the already initalized db
		s.DB = db

		router.Init(s)

		// no need to start echo!
		// see https://github.com/labstack/echo/issues/659

		// if err := s.Start(); err != nil {
		// 	log.Fatalf("Failed to start server: %v", err)
		// }

		// defer func() {
		// 	if err := s.Shutdown(context.Background()); err != nil {
		// 		log.Fatalf("Failed to shutdown server: %v", err)
		// 	}
		// }()

		closure(s)

	})
}
