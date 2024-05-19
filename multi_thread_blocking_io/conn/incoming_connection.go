package conn

import (
	"multi_thread_blocking_io/proto"
	"multi_thread_blocking_io/store"
	"net"
)

// IncomingTCPConnection represents the incoming TCP connection.
type IncomingTCPConnection struct {
	connectionReader      ConnectionReader
	handlersByMessageType map[uint32]Handler
	closeChannel          chan struct{}
}

// NewIncomingTCPConnection creates a new IncomingTCPConnection to handle incoming requests.
func NewIncomingTCPConnection(
	connection net.Conn,
	store *store.InMemoryStore,
) IncomingTCPConnection {
	handlersByMessageType := map[uint32]Handler{
		proto.KeyValueMessageKindPutOrUpdate: NewPutOrUpdateHandler(store),
		proto.KeyValueMessageKindGet:         NewGetHandler(store),
	}
	return IncomingTCPConnection{
		connectionReader:      NewConnectionReader(connection),
		handlersByMessageType: handlersByMessageType,
		closeChannel:          make(chan struct{}),
	}
}

// Handle handles the incoming connection.
func (incomingConnection IncomingTCPConnection) Handle() {
	for {
		select {
		case <-incomingConnection.closeChannel:
			return
		default:
			incomingMessage, err := incomingConnection.connectionReader.AttemptReadOrErrorOut()
			if err != nil {
				return
			}
			switch incomingMessage.Kind {
			case proto.KeyValueMessageKindPutOrUpdate:
				incomingConnection.handlePutOrUpdate(incomingMessage)
			case proto.KeyValueMessageKindGet:
				incomingConnection.handleGet(incomingMessage)
			}
		}
	}
}

// Close closes the IncomingTCPConnection.
func (incomingConnection IncomingTCPConnection) Close() {
	incomingConnection.connectionReader.Close()
	close(incomingConnection.closeChannel)
}

// handlePutOrUpdate handles PutOrUpdate.
func (incomingConnection IncomingTCPConnection) handlePutOrUpdate(message *proto.KeyValueMessage) {
	buffer, err := incomingConnection.handlersByMessageType[message.Kind].Handle(message)
	if err == nil {
		_, _ = incomingConnection.connectionReader.connection.Write(buffer)
	}
}

// handleGet handles Get.
func (incomingConnection IncomingTCPConnection) handleGet(message *proto.KeyValueMessage) {
	buffer, err := incomingConnection.handlersByMessageType[message.Kind].Handle(message)
	if err == nil {
		_, _ = incomingConnection.connectionReader.connection.Write(buffer)
	}
}
