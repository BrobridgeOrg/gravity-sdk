package keyring

import (
	"sync"
)

type Permission struct {
	permissions []string
	mutex       sync.RWMutex
}

func NewPermission() *Permission {
	return &Permission{
		permissions: make([]string, 0),
	}
}

func (p *Permission) Reset() {
	p.permissions = make([]string, 0)
}

func (p *Permission) AddPermissions(permissions []string) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.permissions = append(p.permissions, permissions...)
}

func (p *Permission) GetPermissions() []string {
	return p.permissions
}

func (p *Permission) Check(rule string) bool {

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, entry := range p.permissions {
		if entry == rule {
			return true
		}
	}

	return false
}
