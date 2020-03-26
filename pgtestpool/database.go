package pgtestpool

type Database struct {
	ID           int
	TemplateHash string
	Config       ConnectionConfig
	Closed       bool
	Dirty        bool
	Template     bool
}

func (d *Database) Ready() bool {
	return !d.Closed && !d.Dirty
}

func (d *Database) ReadyForTest() bool {
	return !d.Closed && !d.Dirty && !d.Template
}
