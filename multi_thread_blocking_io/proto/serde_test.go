package proto

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializesAndDeserializesAPutOrUpdateMessage(t *testing.T) {
	message := NewPutOrUpdateKeyValueMessage("DiskType", "SSD")
	buffer, err := message.Serialize()

	assert.Nil(t, err)

	deserializedMessage, err := DeserializeFrom(bytes.NewReader(buffer))

	assert.Nil(t, err)
	assert.Equal(t, "DiskType", deserializedMessage.Key)
	assert.Equal(t, "SSD", deserializedMessage.Value)
	assert.Equal(t, KeyValueMessageKindPutOrUpdate, deserializedMessage.Kind)
}

func TestSerializesAndDeserializesAGetMessage(t *testing.T) {
	message := NewGetValueMessage("DiskType")
	buffer, err := message.Serialize()

	assert.Nil(t, err)

	deserializedMessage, err := DeserializeFrom(bytes.NewReader(buffer))

	assert.Nil(t, err)
	assert.Equal(t, "DiskType", deserializedMessage.Key)
	assert.Equal(t, KeyValueMessageKindGet, deserializedMessage.Kind)
}
