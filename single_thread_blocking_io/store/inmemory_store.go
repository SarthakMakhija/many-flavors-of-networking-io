package store

// InMemoryStore represents a store to hold Key/Value pairs in RAM.
// It is a wrapper over golang's map.
type InMemoryStore struct {
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
	store.valueByKey[key] = value
}

// GetValue gets the value of the given key.
func (store *InMemoryStore) GetValue(key string) (string, bool) {
	value, ok := store.valueByKey[key]
	return value, ok
}
