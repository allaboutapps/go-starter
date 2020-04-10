package pgtestpool

// swagger:response initializeTemplateResponse
// TemplateDatabase description
// in: body
type TemplateDatabase struct {
	// some body description
	// 
	// Required: true
	// Example: Expected type int
	Database `json:"database"`

	// some description for id
	nextTestID    int
	testDatabases []*TestDatabase 
}
