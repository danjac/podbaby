package sql

import (
	"sync"
)

// GenPurse is a literal implementation of a Purse that is programmatically
// generated from SQL file contents within a directory via go generate.
type GenPurse struct {
	mu sync.RWMutex
	files map[string]string
}

// Get takes a filename and returns a query if it is found within the relevant
// map.  If filename is not found, ok will return false.
func (p *GenPurse) Get(filename string) (v string, ok bool) {
	p.mu.RLock()
	v, ok = p.files[filename]
	p.mu.RUnlock()
	return
}
