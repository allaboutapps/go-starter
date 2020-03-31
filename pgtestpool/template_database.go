package pgtestpool

import "sync"

type TemplateDatabase struct {
	Database
	sync.Mutex

	NextTestID    int
	TestDatabases []*TestDatabase
}
