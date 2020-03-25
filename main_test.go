package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"allaboutapps.at/aw/go-mranftl-sample/models"
	"allaboutapps.at/aw/go-mranftl-sample/test"
	"github.com/volatiletech/sqlboiler/boil"
)

func TestInsert(t *testing.T) {

	fmt.Println("Connecting...")

	boil.DebugMode = true

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	_, err = models.Pilots().DeleteAll(context.Background(), db)
	if err != nil {
		t.Error(err)
	}

	pilots, _, _, _ := test.GetFixtures()
	t.Log(pilots)

	tx, err := db.BeginTx(context.TODO(), nil)

	if err != nil {
		t.Error("transaction fail")
	}

	for _, pilot := range pilots {
		pilot.InsertP(context.Background(), db, boil.Infer())
	}

	// Rollback or commit
	err = tx.Commit()

	if err != nil {
		t.Error("transaction commit failed")
	}

}
