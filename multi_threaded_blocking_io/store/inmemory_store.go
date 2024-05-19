package store

import "sync"

// InMemoryStore represents a store to hold Key/Value pairs in RAM.
type InMemoryStore struct {
	lock       sync.RWMutex
	valueByKey map[string]string
}

// NewInMemoryStore creates a new instance if InMemoryStore.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		valueByKey: make(map[string]string),
	}
}

// PutOrUpdate puts or updates the value of the given key.
func (store *InMemoryStore) PutOrUpdate(key, value string) {
	store.lock.Lock()
	store.valueByKey[key] = value
	store.lock.Unlock()
}

// GetValue gets the value of the given key.
func (store *InMemoryStore) GetValue(key string) (string, bool) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	value, ok := store.valueByKey[key]
	return value, ok
}
