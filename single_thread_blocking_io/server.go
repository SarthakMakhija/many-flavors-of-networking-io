package single_thread_blocking_io

import (
	"fmt"
	"log"
	"net"
	_ "net/http/pprof"
	"single_thread_blocking_io/conn"
	"single_thread_blocking_io/store"
)

// TCPServer represents a TCP TCPServer
type TCPServer struct {
	address  string
	listener net.Listener
	store    *store.InMemoryStore
}

// NewTCPServer creates a new instance of TCPServer.
func NewTCPServer(host string, port uint16) (*TCPServer, error) {
	address := fmt.Sprintf("%s:%v", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &TCPServer{
		address:  address,
		listener: listener,
		store:    store.NewInMemoryStore(),
	}, nil
}

// Start starts the server.
// TCPServer implements "Single thread blocking IO" pattern.
// TCPServer:
// - runs a continuous loop in a single goroutine (/main goroutine).
// - a new instance of IncomingTCPConnection is created for every new connection.
// - The incoming TCP connection is handled in the same main goroutine.
// - This pattern involves blocking IO to read from the incoming connection.
func (server *TCPServer) Start() {
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			return
		}
		conn.NewIncomingTCPConnection(connection, server.store).Handle()
	}
}

// Stop stops the server.
func (server *TCPServer) Stop() {
	log.Println("Stopping TCPServer")
	_ = server.listener.Close()
}
