package keyring

import (
	"sync"
)

type Collection struct {
	collections []string
	mutex       sync.RWMutex
}

func NewCollection() *Collection {
	return &Collection{
		collections: make([]string, 0),
	}
}

func (p *Collection) Reset() {
	p.collections = make([]string, 0)
}

func (p *Collection) AddCollections(collections []string) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.collections = append(p.collections, collections...)
}

func (p *Collection) GetCollections() []string {
	return p.collections
}

func (p *Collection) Check(rule string) bool {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Allow all collections
	if len(p.collections) == 0 {
		return true
	}

	for _, entry := range p.collections {
		if entry == rule {
			return true
		}
	}

	return false
}
