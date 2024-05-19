package single_threaded_blocking_io

import (
	"fmt"
	"log"
	"multi_thread_blocking_io/conn"
	"multi_thread_blocking_io/store"
	"net"
	_ "net/http/pprof"
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
func (server *TCPServer) Start() {
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			return
		}
		go conn.NewIncomingTCPConnection(connection, server.store).Handle()
	}
}

// Stop stops the server.
func (server *TCPServer) Stop() {
	log.Println("Stopping TCPServer")
	_ = server.listener.Close()
}
