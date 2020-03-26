package pgtestpool

import (
	"database/sql"
	"testing"
)

// test helpers should never return errors, but are passed the *testing.T instance and fail if needed. It seems to be recommended helper functions are moved to a testing.go file...
// https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859
// https://about.sourcegraph.com/go/advanced-testing-in-go

func disconnectManager(t *testing.T, m *Manager) {
	t.Helper()

	if err := m.Disconnect(true); err != nil {
		t.Logf("received error while disconnecting manager: %v", err)
	}
}

func createTestTables(t *testing.T, template *Database) {
	t.Helper()

	templateDb, err := sql.Open("postgres", template.Config.ConnectionString())
	if err != nil {
		t.Fatalf("failed to connect to template database: %v", err)
	}
	defer templateDb.Close()

	if err := templateDb.Ping(); err != nil {
		t.Fatalf("failed to ping template database: %v", err)
	}

	if _, err := templateDb.Exec(`
		CREATE EXTENSION "uuid-ossp";
		CREATE TABLE pilots (
			id uuid NOT NULL DEFAULT uuid_generate_v4(),
			"name" text NOT NULL,
			created_at timestamptz NOT NULL,
			updated_at timestamptz NULL,
			CONSTRAINT pilot_pkey PRIMARY KEY (id)
		);
		CREATE TABLE jets (
			id uuid NOT NULL DEFAULT uuid_generate_v4(),
			pilot_id uuid NOT NULL,
			age int4 NOT NULL,
			"name" text NOT NULL,
			color text NOT NULL,
			created_at timestamptz NOT NULL,
			updated_at timestamptz NULL,
			CONSTRAINT jet_pkey PRIMARY KEY (id)
		);
		ALTER TABLE jets ADD CONSTRAINT jet_pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id);
	`); err != nil {
		t.Fatalf("failed to create tables in template database: %v", err)
	}

	if _, err := templateDb.Exec(`
		INSERT INTO pilots (id, "name", created_at, updated_at) VALUES
			('744a1a87-5ef7-4309-8814-0f1054751156', 'Mario', '2020-03-23 09:44:00.548', '2020-03-23 09:44:00.548'),
			('20d9d155-2e95-49a2-8889-2ae975a8617e', 'Nick', '2020-03-23 09:44:00.548', '2020-03-23 09:44:00.548');
		INSERT INTO jets (id, pilot_id, age, "name", color, created_at, updated_at) VALUES
			('67d9d0c7-34e5-48b0-9c7d-c6344995353c', '744a1a87-5ef7-4309-8814-0f1054751156', 26, 'F-14B', 'grey', '2020-03-23 09:44:00.000', '2020-03-23 09:44:00.000'),
			('facaf791-21b4-401a-bbac-67079ae4921f', '20d9d155-2e95-49a2-8889-2ae975a8617e', 27, 'F-14B', 'grey/red', '2020-03-23 09:44:00.000', '2020-03-23 09:44:00.000');
	`); err != nil {
		t.Fatalf("failed to insert test data into tables in template database: %v", err)
	}
}

func verifyTestDatabase(t *testing.T, database *Database) {
	t.Helper()

	db, err := sql.Open("postgres", database.Config.ConnectionString())
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	var pilotCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM pilots").Scan(&pilotCount); err != nil {
		t.Fatalf("failed to query pilot test data count: %v", err)
	}

	if pilotCount != 2 {
		t.Errorf("invalid pilot test data count, got %d, want 2", pilotCount)
	}

	var jetCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM jets").Scan(&jetCount); err != nil {
		t.Fatalf("failed to query jet test data count: %v", err)
	}

	if jetCount != 2 {
		t.Errorf("invalid jet test data count, got %d, want 2", jetCount)
	}
}
