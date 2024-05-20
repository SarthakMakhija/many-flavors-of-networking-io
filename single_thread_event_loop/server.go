package single_thread_event_loop

import (
	"log"
	"net"
	"single_thread_eventloop/conn"
	"single_thread_eventloop/proto"
	"single_thread_eventloop/store"
	"syscall"
)

const MaxClients = 10_000

// TCPServer represents an async TCP TCPServer
type TCPServer struct {
	serverFd  int
	eventLoop *EventLoop
}

// NewTCPServer creates a new instance of TCPServer.
func NewTCPServer(host string, port uint16) (*TCPServer, error) {
	//starts the listener on the given port and returns the server file descriptor, if there is no error.
	startListener := func() (int, error) {
		serverFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		if err != nil {
			_ = syscall.Close(serverFd)
			return -1, err
		}
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
	//createEventLoop creates an instance of Event loop.
	createEventLoop := func(serverFd int, store *store.InMemoryStore) (*EventLoop, error) {
		eventLoop, err := NewEventLoop(serverFd, MaxClients, map[uint32]conn.Handler{
			proto.KeyValueMessageKindPutOrUpdate: conn.NewPutOrUpdateHandler(store),
			proto.KeyValueMessageKindGet:         conn.NewGetHandler(store),
		})
		if err != nil {
			return nil, err
		}
		return eventLoop, nil
	}
	//init creates an instance of TCPServer.
	init := func() (*TCPServer, error) {
		serverFd, err := startListener()
		if err != nil {
			return nil, err
		}
		eventLoop, err := createEventLoop(serverFd, store.NewInMemoryStore())
		if err != nil {
			return nil, err
		}
		return &TCPServer{
			serverFd:  serverFd,
			eventLoop: eventLoop,
		}, nil
	}
	return init()
}

// Start starts the server which in turn starts the event loop.
func (server *TCPServer) Start() {
	server.eventLoop.Run()
}

// Stop stops the server.
func (server *TCPServer) Stop() {
	log.Println("Stopping TCPServer")

	server.eventLoop.Stop()
	_ = syscall.Close(server.serverFd)
}
