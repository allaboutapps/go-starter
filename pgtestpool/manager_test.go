package pgtestpool

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"testing"
)

var (
	defaultManagerConfig = ManagerConfig{
		TestDatabaseBaseName: fmt.Sprintf("%s_test", os.Getenv("PSQL_DBNAME")),
		DatabaseConfig: ConnectionConfig{
			Host:     os.Getenv("PSQL_HOST"),
			Port:     5432,
			Username: os.Getenv("PSQL_USER"),
			Password: os.Getenv("PSQL_PASS"),
			Database: os.Getenv("PSQL_DBNAME"),
		},
	}
)

func TestManagerConnectionSuccess(t *testing.T) {
	t.Parallel()

	m := NewManager(defaultManagerConfig)
	if err := m.Connect(); err != nil {
		t.Errorf("manager connection failed: %v", err)
	}

	defer disconnectManager(t, m)

	if !m.Ready() {
		t.Error("manager is not ready")
	}
}

func TestManagerConnectionError(t *testing.T) {
	t.Parallel()

	m := NewManager(ManagerConfig{
		DatabaseConfig: ConnectionConfig{
			Host:     "definitelydoesnotexist",
			Port:     2345,
			Username: "definitelydoesnotexist",
			Password: "definitelydoesnotexist",
			Database: "definitelydoesnotexist",
		},
	})
	if err := m.Connect(); err == nil {
		t.Error("manager connection succeeded")
	}

	if m.Ready() {
		t.Errorf("manager is ready")
	}
}

func TestManagerConnectionAlreadyEstablished(t *testing.T) {
	t.Parallel()

	m := NewManager(defaultManagerConfig)
	if err := m.Connect(); err != nil {
		t.Fatalf("manager connection failed: %v", err)
	}

	defer disconnectManager(t, m)

	if !m.Ready() {
		t.Fatal("manager is not ready")
	}

	if err := m.Connect(); err == nil {
		t.Error("manager connection succeeded although already connected")
	}

	if !m.Ready() {
		t.Error("manager is not ready anymore")
	}
}

func TestManagerReconnectSuccess(t *testing.T) {
	t.Parallel()

	m := NewManager(defaultManagerConfig)
	if err := m.Connect(); err != nil {
		t.Fatalf("manager connection failed: %v", err)
	}

	defer disconnectManager(t, m)

	if !m.Ready() {
		t.Fatal("manager is not ready")
	}

	if err := m.Reconnect(false); err != nil {
		t.Errorf("manager reconnect failed: %v", err)
	}

	if !m.Ready() {
		t.Error("manager is not ready after reconnect")
	}
}

func TestManagerInitTemplateDatabase(t *testing.T) {
	m := NewManager(defaultManagerConfig)
	if err := m.Connect(); err != nil {
		t.Fatalf("manager connection failed: %v", err)
	}

	defer disconnectManager(t, m)

	if !m.Ready() {
		t.Fatal("manager is not ready")
	}

	templateHash := "definitelydoesnotexist"
	template, err := m.InitTemplateDatabase(templateHash, false)
	if err != nil {
		t.Fatalf("failed to initialize template database: %v", err)
	}

	expectedDatabaseName := fmt.Sprintf("%s_template_%s", defaultManagerConfig.DatabaseConfig.Database, templateHash)
	if template.Config.Database != expectedDatabaseName {
		t.Errorf("invalid template database name, got %q, want %q", template.Config.Database, expectedDatabaseName)
	}

	if template.Closed {
		t.Error("template database is flagged as closed")
	}
	if template.Dirty {
		t.Error("template database is flagged as dirty")
	}
	if !template.Template {
		t.Error("template database is not flagged as template")
	}

	if !template.Ready() {
		t.Error("template database is not ready")
	}

	templateHashes := m.GetTemplateDatabaseHashes()
	if len(templateHashes) != 1 {
		t.Fatalf("invalid number of template hashes, got %d, want 1", len(templateHashes))
	}

	if templateHashes[0] != templateHash {
		t.Errorf("invalid template hash, got %q, want %q", templateHashes[0], templateHash)
	}
}

func TestManagerFullTestCycle(t *testing.T) {
	m := NewManager(defaultManagerConfig)
	if err := m.Connect(); err != nil {
		t.Fatalf("manager connection failed: %v", err)
	}

	defer disconnectManager(t, m)

	if !m.Ready() {
		t.Fatal("manager is not ready")
	}

	templateHash := "definitelydoesnotexist"
	template, err := m.InitTemplateDatabase(templateHash, true)
	if err != nil {
		t.Fatalf("failed to initialize template database: %v", err)
	}

	createTestTables(t, template)

	template, err = m.FinalizeTemplateDatabase(templateHash)
	if err != nil {
		t.Fatalf("failed to finalize template database: %v", err)
	}

	if template.Ready() {
		t.Fatal("template database is still ready after finalizing")
	}

	if !template.Closed {
		t.Fatal("template database is not closed after finalizing")
	}

	testDatabaseCount := 10
	if err := m.CreateTestDatabasePool(templateHash, testDatabaseCount); err != nil {
		t.Fatalf("failed to create test database pool: %v", err)
	}

	db, err := sql.Open("postgres", defaultManagerConfig.DatabaseConfig.ConnectionString())
	if err != nil {
		t.Fatalf("failed to open master database connection: %v", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping master database connection: %v", err)
	}

	rows, err := db.Query("SELECT datname FROM pg_database WHERE datname LIKE $1", fmt.Sprintf("%s_%%", defaultManagerConfig.TestDatabaseBaseName))
	if err != nil {
		t.Fatalf("failed to query test database names: %v", err)
	}
	defer rows.Close()

	dbNames := make([]string, 0)
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			t.Fatalf("failed to scan test database row: %v", err)
		}
		dbNames = append(dbNames, dbName)
	}

	sort.Strings(dbNames)

	for i := 0; i < len(dbNames); i++ {
		expectedDBName := fmt.Sprintf("%s_%03d", defaultManagerConfig.TestDatabaseBaseName, i)
		if dbNames[i] != expectedDBName {
			t.Errorf("invalid test database name, got %q, want %q", dbNames[i], expectedDBName)
		}
	}

	testDB, err := m.GetTestDatabaseFromPool(templateHash)
	if err != nil {
		t.Fatalf("failed to get test database from pool: %v", err)
	}

	if !testDB.Dirty {
		t.Error("test database was not flagged as dirty on retrieval")
	}

	verifyTestDatabase(t, testDB)

	testDB2, err := m.GetTestDatabaseFromPool(templateHash)
	if err != nil {
		t.Fatalf("failed to get seconds test database from pool: %v", err)
	}

	if !testDB2.Dirty {
		t.Error("second test database was not flagged as dirty on retrieval")
	}

	verifyTestDatabase(t, testDB2)

	if err := m.ReturnTestDatabaseToPool(testDB, true, false); err != nil {
		t.Fatalf("failed to return test database to pool: %v", err)
	}

	returnedDBName := testDB2.Config.Database

	if err := m.ReturnTestDatabaseToPool(testDB2, false, false); err != nil {
		t.Fatalf("failed to return second test database to pool: %v", err)
	}

	testDB3, err := m.GetTestDatabaseFromPool(templateHash)
	if err != nil {
		t.Fatalf("failed to get third test database from pool: %v", err)
	}

	if !testDB3.Dirty {
		t.Error("third test database was not flagged as dirty on retrieval")
	}

	verifyTestDatabase(t, testDB3)

	if testDB3.Config.Database != returnedDBName {
		t.Errorf("third test database name does not match returned clean second one, got %q, want %q", testDB3.Config.Database, returnedDBName)
	}

	if err := m.ReturnTestDatabaseToPool(testDB3, true, true); err != nil {
		t.Fatalf("failed to return third test database to pool: %v", err)
	}

	var returnedDBExists bool
	if err := m.db.QueryRow("SELECT 1 as exists FROM pg_database WHERE datname = $1", returnedDBName).Scan(&returnedDBExists); err != sql.ErrNoRows || returnedDBExists {
		t.Fatal("third test database was not destroyed upon return")
	}
}
