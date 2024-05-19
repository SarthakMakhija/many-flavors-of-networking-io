package conn

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"non_blocking_busy_waiting/proto"
	store2 "non_blocking_busy_waiting/store"
	"testing"
)

func TestPutAKeyValuePair(t *testing.T) {
	store := store2.NewInMemoryStore()
	handler := NewPutOrUpdateHandler(store)

	putOrUpdateKeyValueMessage := proto.NewPutOrUpdateKeyValueMessage("DiskType", "NVMe")
	handle, err := handler.Handle(putOrUpdateKeyValueMessage)

	assert.Nil(t, err)
	response, _ := proto.DeserializeFrom(bytes.NewReader(handle))

	assert.Equal(t, proto.KeyValueMessageKindPutOrUpdate, response.Kind)
	assert.Equal(t, proto.Status_Ok, response.GetStatus())
}

func TestGetANonExistingKeyValuePair(t *testing.T) {
	store := store2.NewInMemoryStore()
	handler := NewGetHandler(store)

	getValueMessage := proto.NewGetValueMessage("DiskType")
	handle, err := handler.Handle(getValueMessage)

	assert.Nil(t, err)
	response, _ := proto.DeserializeFrom(bytes.NewReader(handle))

	assert.Equal(t, proto.KeyValueMessageKindGetResponse, response.Kind)
	assert.Equal(t, proto.Status_NotOk, response.GetStatus())
}

func TestGetAnExistingKeyValuePair(t *testing.T) {
	store := store2.NewInMemoryStore()
	handler := NewPutOrUpdateHandler(store)

	putOrUpdateKeyValueMessage := proto.NewPutOrUpdateKeyValueMessage("DiskType", "NVMe")
	_, err := handler.Handle(putOrUpdateKeyValueMessage)

	assert.Nil(t, err)

	getValueMessage := proto.NewGetValueMessage("DiskType")
	handle, err := NewGetHandler(store).Handle(getValueMessage)

	assert.Nil(t, err)
	response, _ := proto.DeserializeFrom(bytes.NewReader(handle))

	assert.Equal(t, proto.KeyValueMessageKindGetResponse, response.Kind)
	assert.Equal(t, proto.Status_Ok, response.GetStatus())
	assert.Equal(t, "NVMe", response.GetValue())
}
