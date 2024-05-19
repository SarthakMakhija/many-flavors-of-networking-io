package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPutsAKeyValuePair(t *testing.T) {
	store := NewInMemoryStore()
	store.PutOrUpdate("DiskType", "SSD")

	value, ok := store.GetValue("DiskType")

	assert.True(t, ok)
	assert.Equal(t, "SSD", value)
}

func TestUpdatesTheValueOfAKey(t *testing.T) {
	store := NewInMemoryStore()
	store.PutOrUpdate("DiskType", "SSD")

	store.PutOrUpdate("DiskType", "HDD")
	value, ok := store.GetValue("DiskType")

	assert.True(t, ok)
	assert.Equal(t, "HDD", value)
}

func TestGetsTheValueOfANonExistingKey(t *testing.T) {
	store := NewInMemoryStore()

	value, ok := store.GetValue("DiskType")

	assert.False(t, ok)
	assert.Empty(t, value)
}
