package test

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/test/pgconsumer"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/volatiletech/sqlboiler/boil"
)

var (
	client     *pgconsumer.Client
	hash       string
	migDir, _  = filepath.Abs("../migrations")
	fixFile, _ = filepath.Abs("./fixtures.go")
)

func initIntegres() {
	c, err := pgconsumer.DefaultClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create new pgconsumer client: %v", err)
	}

	client = c
}

func initHash() {
	// migDir, err := filepath.Abs("../migrations")
	// if err != nil {
	// 	log.Fatalf("Failed to get absolute path of migrations directory: %v", err)
	// }
	// fixFile, err := filepath.Abs("./fixtures.go")
	// if err != nil {
	// 	log.Fatalf("Failed to get absolut path of fixtures file: %v", err)
	// }

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
func WithTestServer(closure func(server api.Server)) {
	// TODO
}
