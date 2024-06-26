package conn

import (
	"bytes"
	"errors"
	"io"
	"non_blocking_busy_waiting/proto"
	"syscall"
)

// Client handles an incoming connection (/socket).
type Client struct {
	fd            int
	handlers      map[uint32]Handler
	stopChannel   chan struct{}
	readBuffer    []byte
	currentBuffer *bytes.Buffer
}

// NewClient creates a new instance of the client.
// It reads chunks from the file descriptor and maintains the current buffer.
// currentBuffer denotes the chunk that is read currently.
// The provided file descriptor is set to non-blocking by the caller.
func NewClient(fd int, handlers map[uint32]Handler) *Client {
	return &Client{
		fd:            fd,
		handlers:      handlers,
		stopChannel:   make(chan struct{}),
		readBuffer:    make([]byte, 1024),
		currentBuffer: bytes.NewBuffer([]byte{}),
	}
}

// Run runs the client.
func (client *Client) Run() {
	for {
		select {
		case <-client.stopChannel:
			return
		default:
			keyValueMessage, err := client.read()
			if err != nil {
				return
			}
			if keyValueMessage == nil {
				continue
			}
			if err := client.handle(keyValueMessage); err != nil {
				return
			}
		}
	}
}

// Stop stops the client.
func (client *Client) Stop() {
	close(client.stopChannel)
	_ = syscall.Close(client.fd)
}

// read reads a single proto.KeyValueMessage from the file descriptor.
// The file descriptor is already set to non-blocking, which means syscall.Read(..) will not block.
// However, if there is nothing to be read from the file descriptor, an error would be returned.
// The error would be EAGAIN or EWOULDBLOCK.
// For any error, other than EAGAIN or EWOULDBLOCK, the read method will return.
// read will continue reading till it finds the proto.FooterBytes.
// However, it is possible that syscall.Read(..) does not return the amount of data that is requested.
// In that case, the received data will be stored in client.currentBuffer and the read method will perform poll again.
// When it polls again, it will read further data until proto.FooterBytes are found.
// The combined data represented by the currentBuffer will be deserialized into proto.KeyValueMessage.
func (client *Client) read() (*proto.KeyValueMessage, error) {
	for {
		n, err := syscall.Read(client.fd, client.readBuffer)
		if n <= 0 {
			break
		}
		client.currentBuffer.Write(client.readBuffer[:n])
		if err != nil {
			if err == io.EOF {
				return nil, err
			}
			if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) {
				return nil, nil
			}
			return nil, err
		}
		if bytes.Contains(client.readBuffer, proto.FooterBytes) {
			break
		}
	}
	if client.currentBuffer.Len() > 0 {
		keyValueMessage, err := proto.DeserializeFrom(client.currentBuffer)
		if err != nil {
			return nil, err
		}
		return keyValueMessage, nil
	}
	return nil, nil
}

// handle handles the incoming message.
func (client *Client) handle(keyValueMessage *proto.KeyValueMessage) error {
	buffer, err := client.handlers[keyValueMessage.Kind].Handle(keyValueMessage)
	if err != nil {
		return err
	}
	_, err = client.writeResponse(buffer)
	return err
}

// writeResponse writes the response to the file descriptor.
func (client *Client) writeResponse(buffer []byte) (int, error) {
	return syscall.Write(client.fd, buffer)
}
