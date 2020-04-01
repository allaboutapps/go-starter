package pgtestpool

type TemplateDatabase struct {
	Database

	nextTestID     int
	testDatabases  []*TestDatabase
}
