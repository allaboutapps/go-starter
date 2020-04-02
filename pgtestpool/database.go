package pgtestpool

import "sync"

type Database struct {
	sync.RWMutex `json:"-"`

	TemplateHash string         `json:"templateHash"`
	Config       DatabaseConfig `json:"config"`

	ready bool
	cond  *sync.Cond
}

func (d *Database) Ready() bool {
	d.RLock()
	defer d.RUnlock()

	return d.ready
}

func (d *Database) WaitUntilReady() {
	if d.Ready() {
		return
	}

	if d.cond == nil {
		d.cond = &sync.Cond{L: &sync.Mutex{}}
	}

	d.cond.L.Lock()

	for !d.Ready() {
		d.cond.Wait()
	}

	d.cond.L.Unlock()
}

func (d *Database) FlagAsReady() {
	if d.Ready() {
		return
	}

	d.Lock()
	d.ready = true
	d.Unlock()

	if d.cond != nil {
		d.cond.Broadcast()
	}
}
