package test

// import (
// 	"context"
// 	"database/sql"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"allaboutapps.at/aw/go-mranftl-sample/pgconsumer"
// 	migrate "github.com/rubenv/sql-migrate"
// 	"github.com/volatiletech/sqlboiler/boil"
// )

// func TestMain(m *testing.M) {
// 	migDir, err := filepath.Abs("../migrations")
// 	if err != nil {
// 		log.Fatalf("Failed to get absolute path of migrations directory: %v", err)
// 	}
// 	fixFile, err := filepath.Abs("./fixtures.go")
// 	if err != nil {
// 		log.Fatalf("Failed to get absolut path of fixtures file: %v", err)
// 	}

// 	hash, err := pgconsumer.GetTemplateHash(migDir, fixFile)
// 	if err != nil {
// 		log.Fatalf("Failed to get template hash: %#v", err)
// 	}

// 	c, err := pgconsumer.DefaultClientFromEnv()
// 	if err != nil {
// 		log.Fatalf("Failed to create new pgconsumer client: %v", err)
// 	}

// 	ctx := context.Background()

// 	// ! DISABLE - interferes with pgconsumer paralell running tests!!!
// 	// if err := c.ResetAllTracking(ctx); err != nil {
// 	// 	log.Fatalf("lol: %v", err)
// 	// }

// 	initTemplate := func(db *sql.DB) error {
// 		migrations := &migrate.FileMigrationSource{Dir: migDir}
// 		n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
// 		if err != nil {
// 			return err
// 		}

// 		log.Printf("Applied %d migrations for hash %q", n, hash)

// 		tx, err := db.BeginTx(ctx, nil)
// 		if err != nil {
// 			return err
// 		}

// 		for _, fixture := range fixtures {
// 			if err := fixture.Insert(ctx, db, boil.Infer()); err != nil {
// 				if errr := tx.Rollback(); errr != nil {
// 					return errr
// 				}

// 				return err
// 			}
// 		}

// 		if err := tx.Commit(); err != nil {
// 			return err
// 		}

// 		log.Printf("Inserted %d fixtures for hash %q", len(fixtures), hash)

// 		return nil
// 	}

// 	if err := c.SetupTemplateWithDBClient(ctx, hash, initTemplate); err != nil {
// 		log.Fatalf("Failed to setup template database for hash %q: %v", hash, err)
// 	}

// 	exit := m.Run()

// 	os.Exit(exit)
// }
