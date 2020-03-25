package pgtestpool

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/lib/pq"
)

var (
	ErrManagerNotReady = errors.New("manager is not ready")
)

type Manager struct {
	config         ManagerConfig
	db             *sql.DB
	templates      map[string]*Database
	templateMutex  sync.RWMutex
	databases      map[string][]*Database
	nextDatabaseID map[string]int
	databaseMutex  sync.Mutex
}

func NewManager(config ManagerConfig) *Manager {
	m := &Manager{
		config:         config,
		db:             nil,
		templates:      map[string]*Database{},
		databases:      map[string][]*Database{},
		nextDatabaseID: map[string]int{},
	}

	if len(m.config.TemplateDatabaseBaseName) == 0 {
		m.config.TemplateDatabaseBaseName = fmt.Sprintf("%s_template", m.config.DatabaseConfig.Database)
	}

	if len(m.config.TestDatabaseBaseName) == 0 {
		m.config.TestDatabaseBaseName = m.config.DatabaseConfig.Database
	}

	return m
}

func (m *Manager) Connect() error {
	if m.db != nil {
		return errors.New("manager is already connected")
	}

	db, err := sql.Open("postgres", m.config.DatabaseConfig.ConnectionString())
	if err != nil {
		return errors.Wrap(err, "failed to open manager database connection")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "failed to ping manager database connection")
	}

	m.db = db

	return nil
}

func (m *Manager) Disconnect(ignoreCloseError bool) error {
	if m.db == nil {
		return errors.New("manager is not connected")
	}

	if err := m.db.Close(); err != nil && !ignoreCloseError {
		return errors.Wrap(err, "failed to close database connection")
	}

	m.db = nil

	return nil
}

func (m *Manager) Reconnect(ignoreDisconnectError bool) error {
	if err := m.Disconnect(ignoreDisconnectError); err != nil && !ignoreDisconnectError {
		return errors.Wrap(err, "failed to disconnect manager while reconnecting")
	}

	return m.Connect()
}

func (m *Manager) Ready() bool {
	return m.db != nil
}

func (m *Manager) GetTemplateDatabaseHashes() []string {
	m.templateMutex.RLock()

	hashes := make([]string, 0, len(m.templates))

	for hash := range m.templates {
		hashes = append(hashes, hash)
	}

	m.templateMutex.RUnlock()

	return hashes
}

func (m *Manager) InitTemplateDatabase(hash string, recreate bool) (*Database, error) {
	if !m.Ready() {
		return nil, ErrManagerNotReady
	}

	m.templateMutex.Lock()
	defer m.templateMutex.Unlock()

	if _, ok := m.templates[hash]; ok {
		return nil, errors.New("cannot to initialize template database, hash already exists")
	}

	templateDatabaseName := fmt.Sprintf("%s_%s", m.config.TemplateDatabaseBaseName, hash)

	var exists bool
	err := m.db.QueryRow("SELECT 1 as exists FROM pg_database WHERE datname = $1", templateDatabaseName).Scan(&exists)
	if err == sql.ErrNoRows {
		exists = false
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query whether template database exists")
	}

	if exists && recreate {
		if _, err := m.db.Exec(fmt.Sprintf("DROP DATABASE %s", pq.QuoteIdentifier(templateDatabaseName))); err != nil {
			return nil, errors.Wrap(err, "failed to drop template database")
		}
	}

	if !exists || recreate {
		if _, err := m.db.Exec(fmt.Sprintf("CREATE DATABASE %s WITH OWNER %s TEMPLATE \"template0\"", pq.QuoteIdentifier(templateDatabaseName), pq.QuoteIdentifier(m.config.DatabaseConfig.Username))); err != nil {
			return nil, errors.Wrap(err, "failed to create template database")
		}
	}

	db := &Database{
		Config: ConnectionConfig{
			Host:     m.config.DatabaseConfig.Host,
			Port:     m.config.DatabaseConfig.Port,
			Username: m.config.DatabaseConfig.Username,
			Password: m.config.DatabaseConfig.Password,
			Database: templateDatabaseName,
		},
		Closed:   false,
		Dirty:    false,
		Template: true,
	}

	m.templates[hash] = db

	return db, nil
}

func (m *Manager) FinalizeTemplateDatabase(hash string) (*Database, error) {
	if !m.Ready() {
		return nil, ErrManagerNotReady
	}

	m.templateMutex.Lock()
	defer m.templateMutex.Unlock()

	template, ok := m.templates[hash]
	if !ok {
		return nil, errors.New("cannot finalize template database, hash does not exist")
	}

	if template.Closed {
		return nil, errors.New("cannot finalize template database, already flagged as closed")
	}

	template.Closed = true

	return template, nil
}

func (m *Manager) CreateTestDatabasePool(hash string, count int) error {
	if !m.Ready() {
		return ErrManagerNotReady
	}

	m.templateMutex.RLock()
	template, ok := m.templates[hash]
	m.templateMutex.RUnlock()

	if !ok {
		return errors.New("cannot create test database pool, template hash does not exist")
	}

	if !template.Closed {
		return errors.New("cannot create test database pool, template database has not been finalized yet")
	}

	m.databaseMutex.Lock()
	defer m.databaseMutex.Unlock()

	if _, ok := m.databases[hash]; !ok {
		m.databases[hash] = make([]*Database, 0)
	}

	if _, ok := m.nextDatabaseID[hash]; !ok {
		m.nextDatabaseID[hash] = 0
	}

	for i := 0; i < count; i++ {
		db, err := m.createTestDatabase(hash, template.Config.Database, m.nextDatabaseID[hash])
		if err != nil {
			return errors.Wrap(err, "failed to create test database for pool")
		}

		m.databases[hash] = append(m.databases[hash], db)
		m.nextDatabaseID[hash]++
	}

	return nil
}

func (m *Manager) GetTestDatabaseFromPool(hash string) (*Database, error) {
	if !m.Ready() {
		return nil, ErrManagerNotReady
	}

	m.templateMutex.RLock()
	template, ok := m.templates[hash]
	m.templateMutex.RUnlock()

	if !ok {
		return nil, errors.New("cannot get test database from pool, template hash does not exist")
	}

	m.databaseMutex.Lock()
	defer m.databaseMutex.Unlock()

	if _, ok := m.databases[hash]; !ok {
		return nil, errors.New("no test database pool created for template hash")
	}

	if _, ok := m.nextDatabaseID[hash]; !ok {
		return nil, errors.New("no next database ID available for template hash")
	}

	var testDB *Database
	for _, db := range m.databases[hash] {
		if db.ReadyForTest() {
			testDB = db
			break
		}
	}

	if testDB == nil {
		db, err := m.createTestDatabase(hash, template.Config.Database, m.nextDatabaseID[hash])
		if err != nil {
			return nil, errors.Wrap(err, "no ready test database available, failed to create fresh one")
		}

		m.databases[hash] = append(m.databases[hash], db)
		m.nextDatabaseID[hash]++

		testDB = db
	}

	testDB.Dirty = true

	newDB, err := m.createTestDatabase(hash, template.Config.Database, m.nextDatabaseID[hash])
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new test database after retrieving one from pool")
	}

	m.databases[hash] = append(m.databases[hash], newDB)
	m.nextDatabaseID[hash]++

	return testDB, nil
}

func (m *Manager) ReturnTestDatabaseToPool(db *Database, dirty bool, destroy bool) error {
	if !m.Ready() {
		return ErrManagerNotReady
	}

	m.databaseMutex.Lock()
	defer m.databaseMutex.Unlock()

	if _, ok := m.databases[db.TemplateHash]; !ok {
		return errors.New("no pool created for template hash, cannot return test database")
	}

	if destroy {
		idx := -1
		for i, testDB := range m.databases[db.TemplateHash] {
			if testDB.ID == db.ID {
				if err := m.destroyTestDatabase(testDB); err != nil {
					return errors.Wrap(err, "failed to destroy test database after returning to pool")
				}

				idx = i
				break
			}
		}

		if idx < 0 {
			return errors.New("test database not found for template hash, cannot destroy")
		}

		// Delete while preserving order without causing memory leaks due to pointers, according to: https://github.com/golang/go/wiki/SliceTricks
		if idx < len(m.databases[db.TemplateHash])-1 {
			copy(m.databases[db.TemplateHash][idx:], m.databases[db.TemplateHash][idx+1:])
		}
		m.databases[db.TemplateHash][len(m.databases[db.TemplateHash])-1] = nil
		m.databases[db.TemplateHash] = m.databases[db.TemplateHash][:len(m.databases[db.TemplateHash])-1]

		return nil
	}

	for _, testDB := range m.databases[db.TemplateHash] {
		if testDB.ID == db.ID {
			testDB.Dirty = dirty
			testDB.Closed = dirty
		}
	}

	return nil
}

func (m *Manager) createTestDatabase(hash string, templateDatabaseName string, id int) (*Database, error) {
	if !m.Ready() {
		return nil, ErrManagerNotReady
	}

	testDatabaseName := fmt.Sprintf("%s_%03d", m.config.TestDatabaseBaseName, id)

	if _, err := m.db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", pq.QuoteIdentifier(testDatabaseName))); err != nil {
		return nil, errors.Wrap(err, "failed to drop test database")
	}

	if _, err := m.db.Exec(fmt.Sprintf("CREATE DATABASE %s WITH OWNER %s TEMPLATE %s", pq.QuoteIdentifier(testDatabaseName), pq.QuoteIdentifier(m.config.DatabaseConfig.Username), pq.QuoteIdentifier(templateDatabaseName))); err != nil {
		return nil, errors.Wrap(err, "failed to create test database")
	}

	db := &Database{
		ID: id,
		TemplateHash: hash,
		Config: ConnectionConfig{
			Host:     m.config.DatabaseConfig.Host,
			Port:     m.config.DatabaseConfig.Port,
			Username: m.config.DatabaseConfig.Username,
			Password: m.config.DatabaseConfig.Password,
			Database: testDatabaseName,
		},
		Closed:   false,
		Dirty:    false,
		Template: false,
	}

	return db, nil
}

func (m *Manager) destroyTestDatabase(db *Database) error {
	if !m.Ready() {
		return ErrManagerNotReady
	}

	if _, err := m.db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", pq.QuoteIdentifier(db.Config.Database))); err != nil {
		return errors.Wrap(err, "failed to destroy test database")
	}

	return nil
}
