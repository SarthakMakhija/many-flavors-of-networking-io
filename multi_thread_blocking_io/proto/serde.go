package proto

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"io"
	"unsafe"
)

const ReservedHeaderLength = int(unsafe.Sizeof(uint32(0)))

var (
	FooterBytes  = []byte{'@', 'E', 'O', 'F', '@'}
	FooterLength = len(FooterBytes)
)

const (
	KeyValueMessageKindGet         = uint32(1)
	KeyValueMessageKindGetResponse = uint32(2)
	KeyValueMessageKindPutOrUpdate = uint32(3)
)

// NewPutOrUpdateKeyValueMessage creates a new instance of KeyValueMessage with kind as PutOrUpdate.
func NewPutOrUpdateKeyValueMessage(key, value string) *KeyValueMessage {
	return &KeyValueMessage{
		Key:   key,
		Value: value,
		Kind:  KeyValueMessageKindPutOrUpdate,
	}
}

// NewGetValueMessage a new instance of KeyValueMessage with kind as Get.
func NewGetValueMessage(key string) *KeyValueMessage {
	return &KeyValueMessage{
		Key:  key,
		Kind: KeyValueMessageKindGet,
	}
}

// NewPutOrUpdateKeyValueSuccessfulResponseMessage creates a new instance of KeyValueMessage with kind as PutOrUpdate.
func NewPutOrUpdateKeyValueSuccessfulResponseMessage() *KeyValueMessage {
	return &KeyValueMessage{
		Kind:   KeyValueMessageKindPutOrUpdate,
		Status: Status_Ok,
	}
}

// NewGetValueSuccessfulResponseMessage creates a new instance of KeyValueMessage with kind as GetResponse.
func NewGetValueSuccessfulResponseMessage(key string, value string) *KeyValueMessage {
	return &KeyValueMessage{
		Key:    key,
		Value:  value,
		Kind:   KeyValueMessageKindGetResponse,
		Status: Status_Ok,
	}
}

// NewGetValueUnsuccessfulResponseMessage creates a new instance of KeyValueMessage with kind as GetResponse.
func NewGetValueUnsuccessfulResponseMessage(key string) *KeyValueMessage {
	return &KeyValueMessage{
		Key:    key,
		Kind:   KeyValueMessageKindGetResponse,
		Status: Status_NotOk,
	}
}

// Serialize serializes the KeyValueMessage in bytes.
// KeyValueMessage is serialized in the following format:
// 4 bytes to denote size -> message.serialize() -> FooterBytes
func (message *KeyValueMessage) Serialize() ([]byte, error) {
	payload, err := message.serialize()
	if err != nil {
		return nil, err
	}

	body := make([]byte, ReservedHeaderLength+len(payload)+FooterLength)
	binary.LittleEndian.PutUint32(body, uint32(len(payload))+uint32(FooterLength))
	copy(body[ReservedHeaderLength:], payload)
	copy(body[ReservedHeaderLength+len(payload):], FooterBytes)

	return body, nil
}

// DeserializeFrom deserializes the reader into KeyValueMessage.
// Usually the incoming connection is passed as a reader.
func DeserializeFrom(reader io.Reader) (*KeyValueMessage, error) {
	headerBytes := make([]byte, ReservedHeaderLength)
	_, err := reader.Read(headerBytes)
	if err != nil {
		return nil, err
	}

	bodyWithFooter := make([]byte, binary.LittleEndian.Uint32(headerBytes))
	_, err = reader.Read(bodyWithFooter)
	if err != nil {
		return nil, err
	}

	message := &KeyValueMessage{}
	err = proto.Unmarshal(bodyWithFooter[:len(bodyWithFooter)-FooterLength], message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// serialize uses proto.Marshal to serialize KeyValueMessage.
func (message *KeyValueMessage) serialize() ([]byte, error) {
	buffer, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
