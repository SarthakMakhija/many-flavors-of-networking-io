package conn

import (
	proto2 "many-flavors-of-nwing-io/multi_threaded_blocking_io/proto"
	"many-flavors-of-nwing-io/multi_threaded_blocking_io/store"
)

// Handler handles the incoming requests.
type Handler interface {
	Handle(message *proto2.KeyValueMessage) ([]byte, error)
}

// PutOrUpdateHandler handles the PutOrUpdate request.
type PutOrUpdateHandler struct {
	store *store.InMemoryStore
}

// NewPutOrUpdateHandler creates a new instance of PutOrUpdateHandler.
func NewPutOrUpdateHandler(store *store.InMemoryStore) Handler {
	return PutOrUpdateHandler{
		store: store,
	}
}

func (handler PutOrUpdateHandler) Handle(message *proto2.KeyValueMessage) ([]byte, error) {
	handler.store.PutOrUpdate(message.Key, message.Value)
	return proto2.NewPutOrUpdateKeyValueSuccessfulResponseMessage().Serialize()
}

// GetHandler handles the Get request.
type GetHandler struct {
	store *store.InMemoryStore
}

// NewGetHandler creates a new instance of GetHandler.
func NewGetHandler(store *store.InMemoryStore) Handler {
	return GetHandler{
		store: store,
	}
}

func (handler GetHandler) Handle(message *proto2.KeyValueMessage) ([]byte, error) {
	value, ok := handler.store.GetValue(message.Key)
	var buffer []byte
	var err error

	if !ok {
		buffer, err = proto2.NewGetValueUnsuccessfulResponseMessage(message.Key).Serialize()
	} else {
		buffer, err = proto2.NewGetValueSuccessfulResponseMessage(message.Key, value).Serialize()
	}
	return buffer, err
}
