package pgtestpool

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"sync"

	"github.com/friendsofgo/errors"
	"github.com/lib/pq"
)

var (
	ErrManagerNotReady            = errors.New("manager not ready")
	ErrTemplateAlreadyInitialized = errors.New("template is already initialized")
	ErrTemplateDoesNotExist       = errors.New("template does not exist")
)

const (
	prefixTemplateDatabase   = "template"
	prefixTestDatabase       = "test"
	templateDatabaseTemplate = "template0"
)

type Manager struct {
	config        ManagerConfig
	db            *sql.DB
	templates     map[string]*TemplateDatabase
	templateMutex sync.RWMutex
	wg            sync.WaitGroup
}

func NewManager(config ManagerConfig) *Manager {
	m := &Manager{
		config:    config,
		db:        nil,
		templates: map[string]*TemplateDatabase{},
		wg:        sync.WaitGroup{},
	}

	if len(m.config.TestDatabaseOwner) == 0 {
		m.config.TestDatabaseOwner = m.config.ManagerDatabaseConfig.Username
	}

	if len(m.config.TestDatabaseOwnerPassword) == 0 {
		m.config.TestDatabaseOwnerPassword = m.config.ManagerDatabaseConfig.Password
	}

	return m
}

func DefaultManagerFromEnv() *Manager {
	return NewManager(DefaultManagerConfigFromEnv())
}

func (m *Manager) Connect(ctx context.Context) error {
	if m.db != nil {
		return errors.New("manager is already connected")
	}

	db, err := sql.Open("postgres", m.config.ManagerDatabaseConfig.ConnectionString())
	if err != nil {
		return errors.Wrap(err, "failed to open manager database connection")
	}

	if err := db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "failed to ping manager database connection")
	}

	m.db = db

	return nil
}

func (m *Manager) Disconnect(ctx context.Context, ignoreCloseError bool) error {
	if m.db == nil {
		return errors.New("manager is not connected")
	}

	c := make(chan struct{})
	go func() {
		defer close(c)
		m.wg.Wait()
	}()

	select {
	case <-c:
	case <-ctx.Done():
	}

	if err := m.db.Close(); err != nil && !ignoreCloseError {
		return errors.Wrap(err, "failed to close database connection")
	}

	m.db = nil

	return nil
}

func (m *Manager) Reconnect(ctx context.Context, ignoreDisconnectError bool) error {
	if err := m.Disconnect(ctx, ignoreDisconnectError); err != nil && !ignoreDisconnectError {
		return errors.Wrap(err, "failed to disconnect manager while reconnecting")
	}

	return m.Connect(ctx)
}

func (m *Manager) Ready() bool {
	return m.db != nil
}

func (m *Manager) Initialize(ctx context.Context) error {
	if !m.Ready() {
		if err := m.Connect(ctx); err != nil {
			return errors.Wrap(err, "failed to connect manager while initializing")
		}
	}

	rows, err := m.db.QueryContext(ctx, "SELECT datname FROM pg_database WHERE datname LIKE $1", fmt.Sprintf("%s_%s_%%", m.config.DatabasePrefix, prefixTestDatabase))
	if err != nil {
		return errors.Wrap(err, "failed to query for existing test databases")
	}
	defer rows.Close()

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return errors.Wrap(err, "failed to scan existing test database row")
		}

		if _, err := m.db.Exec(fmt.Sprintf("DROP DATABASE %s", pq.QuoteIdentifier(dbName))); err != nil {
			return errors.Wrapf(err, "failed to drop existing test database %q", dbName)
		}
	}

	return nil
}

func (m *Manager) InitializeTemplateDatabase(ctx context.Context, hash string) (*TemplateDatabase, error) {
	if !m.Ready() {
		return nil, ErrManagerNotReady
	}

	m.templateMutex.RLock()
	template, ok := m.templates[hash]
	m.templateMutex.RUnlock()

	if ok {
		if template.Ready() {
			return template, nil
		}

		return nil, ErrTemplateAlreadyInitialized
	}

	dbName := fmt.Sprintf("%s_%s_%s", m.config.DatabasePrefix, prefixTemplateDatabase, hash)
	template = &TemplateDatabase{
		Database: Database{
			TemplateHash: hash,
			Config: DatabaseConfig{
				Host:     m.config.ManagerDatabaseConfig.Host,
				Port:     m.config.ManagerDatabaseConfig.Port,
				Username: m.config.ManagerDatabaseConfig.Username,
				Password: m.config.ManagerDatabaseConfig.Password,
				Database: dbName,
			},
			ready: false,
			mutex: &sync.RWMutex{},
		},
		NextTestID:    0,
		TestDatabases: make([]*TestDatabase, 0),
	}

	m.templateMutex.Lock()
	defer m.templateMutex.Unlock()

	m.templates[hash] = template

	if err := m.dropAndCreateDatabase(ctx, dbName, m.config.ManagerDatabaseConfig.Username, templateDatabaseTemplate); err != nil {
		m.templates[hash] = nil

		return nil, errors.Wrapf(err, "failed to drop and create template database %q", dbName)
	}

	return template, nil
}

func (m *Manager) FinalizeTemplateDatabase(ctx context.Context, hash string) (*TemplateDatabase, error) {
	if !m.Ready() {
		return nil, ErrManagerNotReady
	}

	m.templateMutex.Lock()
	defer m.templateMutex.Unlock()

	template, ok := m.templates[hash]
	if !ok {
		dbName := fmt.Sprintf("%s_%s_%s", m.config.DatabasePrefix, prefixTemplateDatabase, hash)
		exists, err := m.checkDatabaseExists(ctx, dbName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to check whether template database %q exists while finalizing", dbName)
		}

		if !exists {
			return nil, errors.Errorf("failed to finalize template database, hash %q does not exist", hash)
		}

		template = &TemplateDatabase{
			Database: Database{
				TemplateHash: hash,
				Config: DatabaseConfig{
					Host:     m.config.ManagerDatabaseConfig.Host,
					Port:     m.config.ManagerDatabaseConfig.Port,
					Username: m.config.ManagerDatabaseConfig.Username,
					Password: m.config.ManagerDatabaseConfig.Password,
					Database: dbName,
				},
				ready: false,
				mutex: &sync.RWMutex{},
			},
			NextTestID:    0,
			TestDatabases: make([]*TestDatabase, 0),
		}

		m.templates[hash] = template
	}

	template.FlagAsReady()

	return template, nil
}

func (m *Manager) GetTestDatabase(ctx context.Context, hash string) (*TestDatabase, error) {
	if !m.Ready() {
		return nil, ErrManagerNotReady
	}

	m.templateMutex.RLock()
	template, ok := m.templates[hash]
	m.templateMutex.RUnlock()

	if !ok {
		return nil, ErrTemplateDoesNotExist
	}

	template.WaitUntilReady()

	template.Lock()
	defer template.Unlock()

	var testDB *TestDatabase
	for _, db := range template.TestDatabases {
		if db.ReadyForTest() {
			testDB = db
			break
		}
	}

	if testDB == nil {
		dbName := fmt.Sprintf("%s_%s_%s_%d", m.config.DatabasePrefix, prefixTestDatabase, template.TemplateHash, template.NextTestID)

		if err := m.createDatabase(ctx, dbName, m.config.TestDatabaseOwner, template.Config.Database); err != nil {
			return nil, errors.Wrapf(err, "failed to create test database %q while retrieving one for hash %q", dbName, hash)
		}

		testDB = &TestDatabase{
			Database: Database{
				TemplateHash: template.TemplateHash,
				Config: DatabaseConfig{
					Host:     m.config.ManagerDatabaseConfig.Host,
					Port:     m.config.ManagerDatabaseConfig.Port,
					Username: m.config.TestDatabaseOwner,
					Password: m.config.TestDatabaseOwnerPassword,
					Database: dbName,
				},
				ready: true,
				mutex: &sync.RWMutex{},
			},
			ID:    template.NextTestID,
			dirty: false,
		}

		template.TestDatabases = append(template.TestDatabases, testDB)
		template.NextTestID++
	}

	testDB.FlagAsDirty()

	m.wg.Add(1)
	go m.addTestDatabaseInBackground(template)

	return testDB, nil
}

func (m *Manager) ReturnTestDatabase(ctx context.Context, hash string, id int) error {
	if !m.Ready() {
		return ErrManagerNotReady
	}

	m.templateMutex.RLock()
	template, ok := m.templates[hash]
	m.templateMutex.RUnlock()

	if !ok {
		return ErrTemplateDoesNotExist
	}

	template.WaitUntilReady()

	template.Lock()
	defer template.Unlock()

	found := false
	for _, db := range template.TestDatabases {
		if db.ID == id {
			found = true
			db.FlagAsClean()
			break
		}
	}

	if !found {
		dbName := fmt.Sprintf("%s_%s_%s_%03d", m.config.DatabasePrefix, prefixTestDatabase, hash, id)
		exists, err := m.checkDatabaseExists(ctx, dbName)
		if err != nil {
			return errors.Wrapf(err, "failed to check whether test database %q exists while returning", dbName)
		}

		if !exists {
			return errors.Errorf("failed to return test database %d for template %q", id, hash)
		}

		db := &TestDatabase{
			Database: Database{
				TemplateHash: hash,
				Config: DatabaseConfig{
					Host:     m.config.ManagerDatabaseConfig.Host,
					Port:     m.config.ManagerDatabaseConfig.Port,
					Username: m.config.TestDatabaseOwner,
					Password: m.config.TestDatabaseOwnerPassword,
					Database: dbName,
				},
				ready: true,
				mutex: &sync.RWMutex{},
			},
			ID: id,
			dirty: false,
		}

		template.TestDatabases = append(template.TestDatabases, db)
		sort.Sort(ByID(template.TestDatabases))
	}

	return nil
}

func (m *Manager) checkDatabaseExists(ctx context.Context, dbName string) (bool, error) {
	var exists bool
	if err := m.db.QueryRowContext(ctx, "SELECT 1 AS exists FROM pg_database WHERE datname = $1", dbName).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, errors.Wrapf(err, "failed to check whether database %q exists", dbName)
	}

	return exists, nil
}

func (m *Manager) createDatabase(ctx context.Context, dbName string, owner string, template string) error {
	if _, err := m.db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s WITH OWNER %s TEMPLATE %s", pq.QuoteIdentifier(dbName), pq.QuoteIdentifier(owner), pq.QuoteIdentifier(template))); err != nil {
		return errors.Wrapf(err, "failed to create database %q", dbName)
	}

	return nil
}

func (m *Manager) addTestDatabaseInBackground(template *TemplateDatabase) {
	defer m.wg.Done()

	// TODO avoid code duplication somehow without running into mutex issues?
	template.Lock()
	defer template.Unlock()

	dbName := fmt.Sprintf("%s_%s_%s_%03d", m.config.DatabasePrefix, prefixTestDatabase, template.TemplateHash, template.NextTestID)

	if err := m.createDatabase(context.Background(), dbName, m.config.TestDatabaseOwner, template.Config.Database); err != nil {
		// TODO log error somewhere instead of silently swallowing it?
		return
	}

	testDB := &TestDatabase{
		Database: Database{
			TemplateHash: template.TemplateHash,
			Config: DatabaseConfig{
				Host:     m.config.ManagerDatabaseConfig.Host,
				Port:     m.config.ManagerDatabaseConfig.Port,
				Username: m.config.TestDatabaseOwner,
				Password: m.config.TestDatabaseOwnerPassword,
				Database: dbName,
			},
			ready: true,
			mutex: &sync.RWMutex{},
		},
		ID:    template.NextTestID,
		dirty: false,
	}

	template.TestDatabases = append(template.TestDatabases, testDB)
	template.NextTestID++
}

func (m *Manager) dropDatabase(ctx context.Context, dbName string) error {
	if _, err := m.db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", pq.QuoteIdentifier(dbName))); err != nil {
		return errors.Wrapf(err, "failed to drop database %q", dbName)
	}

	return nil
}

func (m *Manager) dropAndCreateDatabase(ctx context.Context, dbName string, owner string, template string) error {
	if err := m.dropDatabase(ctx, dbName); err != nil {
		return errors.Wrapf(err, "failed to drop database %q before recreating", dbName)
	}

	return m.createDatabase(ctx, dbName, owner, template)
}
