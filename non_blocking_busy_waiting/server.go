package non_blocking_busy_waiting

import (
	"errors"
	"log"
	"net"
	"non_blocking_busy_waiting/conn"
	"non_blocking_busy_waiting/proto"
	store2 "non_blocking_busy_waiting/store"
	"syscall"
)

const MaxClients = 10_000

// TCPServer represents a non-blocking busy-waiting TCP TCPServer
type TCPServer struct {
	serverFd    int
	handlers    map[uint32]conn.Handler
	stopChannel chan struct{}
}

// NewTCPServer creates a new instance of TCPServer.
func NewTCPServer(host string, port uint16) (*TCPServer, error) {
	//starts the listener on the given port and returns the server file descriptor, if there is no error.
	startListener := func() (int, error) {
		// syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0) creates an IPv4 (AF_INET), bidirectional (SOCK_STREAM), TCP (0) socket.
		serverFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		if err != nil {
			_ = syscall.Close(serverFd)
			return -1, err
		}
		// SetNonblock sets the server file descriptor non-blocking. This means the file descriptor can be polled.
		// A non-blocking file descriptor does not block on IO operations and can be polled.
		if err = syscall.SetNonblock(serverFd, true); err != nil {
			_ = syscall.Close(serverFd)
			return -1, err
		}

		ip4 := net.ParseIP(host)
		if err = syscall.Bind(serverFd, &syscall.SockaddrInet4{
			Port: int(port),
			Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
		}); err != nil {
			_ = syscall.Close(serverFd)
			return -1, err
		}
		if err = syscall.Listen(serverFd, MaxClients); err != nil {
			return -1, err
		}
		return serverFd, nil
	}
	serverFd, err := startListener()
	if err != nil {
		return nil, err
	}

	store := store2.NewInMemoryStore()
	return &TCPServer{
		serverFd: serverFd,
		handlers: map[uint32]conn.Handler{
			proto.KeyValueMessageKindPutOrUpdate: conn.NewPutOrUpdateHandler(store),
			proto.KeyValueMessageKindGet:         conn.NewGetHandler(store),
		},
		stopChannel: make(chan struct{}),
	}, nil
}

// Start starts the server.
// TCPServer implements "Non-Blocking with Busy-Wait" pattern.
// TCPServer:
// - runs a continuous loop in a single goroutine (/main goroutine).
// - serverFd is already marked non-blocking, this means any IO operations on this file descriptor will not block. However, the file descriptor can be polled.
// - an incoming connection is represented by its file descriptor "connectionFd".
// - connectionFd is also marked non-blocking.
// - a new client is created (for the incoming connectionFd) which handles the connection.
// - all the IO operations are non-blocking.
// This server handles only one client at a time.
func (server *TCPServer) Start() {
	for {
		select {
		case <-server.stopChannel:
			return
		default:
			connectionFd, _, err := syscall.Accept(server.serverFd)
			if err != nil {
				if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) {
					continue
				}
				return
			}
			_ = syscall.SetNonblock(connectionFd, true)
			conn.NewClient(connectionFd, server.handlers).Run()
		}
	}
}

// Stop stops the server.
func (server *TCPServer) Stop() {
	log.Println("Stopping TCPServer")
	_ = syscall.Close(server.serverFd)
	close(server.stopChannel)
}
