package pgtestpool

import "sync"

type Database struct {
	TemplateHash string
	Config       DatabaseConfig
	ready        bool
	mutex        *sync.RWMutex
	cond         *sync.Cond
}

func (d *Database) Ready() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

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

	d.mutex.Lock()
	d.ready = true
	d.mutex.Unlock()

	if d.cond != nil {
		d.cond.Broadcast()
	}
}
