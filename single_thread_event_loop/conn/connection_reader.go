package conn

import (
	"bufio"
	"errors"
	"net"
	"single_thread_eventloop/proto"
	"time"
)

const maxTimeoutErrorsTolerable = 10

// ConnectionReader represents an abstraction to read from the connection.
type ConnectionReader struct {
	connection     net.Conn
	closeChannel   chan struct{}
	bufferedReader *bufio.Reader
	netError       net.Error
}

// NewConnectionReader creates a new instance of ConnectionReader.
func NewConnectionReader(connection net.Conn) ConnectionReader {
	return ConnectionReader{
		connection:     connection,
		bufferedReader: bufio.NewReader(connection),
		closeChannel:   make(chan struct{}),
	}
}

// AttemptReadOrErrorOut attempts to read from the incoming TCP connection.
//
// It is only used for testing in single_thread_event_loop module.
//
// It runs an infinite loop to read a single message from the incoming connection.
//
// proto.DeserializeFrom() reads from the connection using "blocking IO" and returns either a message or an error.
// The method tolerates network timeout errors.
//
// This method also sets ReadDeadline for future Read calls and any currently-blocked Read call.
func (connectionReader ConnectionReader) AttemptReadOrErrorOut() (*proto.KeyValueMessage, error) {
	totalTimeoutsErrors := 0
	for {
		select {
		case <-connectionReader.closeChannel:
			return nil, errors.New("ConnectionReader is closed")
		default:
			_ = connectionReader.connection.SetReadDeadline(time.Now().Add(20 * time.Millisecond))

			message, err := proto.DeserializeFrom(connectionReader.bufferedReader)
			if err != nil {
				if errors.As(err, &connectionReader.netError) && connectionReader.netError.Timeout() {
					totalTimeoutsErrors += 1
				}
				if totalTimeoutsErrors <= maxTimeoutErrorsTolerable {
					continue
				}
				return nil, err
			}
			return message, nil
		}
	}
}

// Close closes the ConnectionReader.
func (connectionReader ConnectionReader) Close() {
	close(connectionReader.closeChannel)
}
