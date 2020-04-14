package pgconsumer

type TestDatabase struct {
	Database `json:"database"`

	ID int `json:"id"`
}
