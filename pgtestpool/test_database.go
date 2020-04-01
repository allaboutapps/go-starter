package pgtestpool

type TestDatabase struct {
	Database

	ID int

	dirty bool
}

func (t *TestDatabase) Dirty() bool {
	t.RLock()
	defer t.RUnlock()

	return t.dirty
}

func (t *TestDatabase) FlagAsDirty() {
	t.Lock()
	defer t.Unlock()

	t.dirty = true
}

func (t *TestDatabase) FlagAsClean() {
	t.Lock()
	defer t.Unlock()

	t.dirty = false
}

func (t *TestDatabase) ResetStatus() {
	t.Lock()
	t.dirty = false
	t.ready = true
	t.Unlock()

	if t.cond != nil {
		t.cond.Broadcast()
	}
}

func (t *TestDatabase) ReadyForTest() bool {
	return t.Ready() && !t.Dirty()
}

type ByID []*TestDatabase

func (i ByID) Len() int           { return len(i) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }
