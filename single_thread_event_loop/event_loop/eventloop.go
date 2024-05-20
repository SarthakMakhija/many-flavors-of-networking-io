package event_loop

import (
	"single_thread_eventloop/conn"
	"syscall"
)

// EventLoop represents a single goroutine event loop.
type EventLoop struct {
	serverFd       int
	kQueue         *KQueue
	clients        map[int]*Client
	clientHandlers map[uint32]conn.Handler
	stopChannel    chan struct{}
}

// NewEventLoop creates a new instance of EventLoop.
// It also subscribes using the EVFILT_READ filter on the server file descriptor.
func NewEventLoop(serverFd int, maxClients int, clientHandlers map[uint32]conn.Handler) (*EventLoop, error) {
	// NewKQueue creates a new kernel KQueue data structure to hold various events on the subscribed file descriptor.
	kQueue, err := NewKQueue(maxClients)
	if err != nil {
		return nil, err
	}
	eventLoop := &EventLoop{
		serverFd:       serverFd,
		kQueue:         kQueue,
		clients:        make(map[int]*Client),
		clientHandlers: clientHandlers,
		stopChannel:    make(chan struct{}),
	}
	// subscribes to the given server file descriptor using EVFILT_READ and EV_ADD flag.
	// This means an event will be added to the kernel KQueue when the server file descriptor is ready to be read
	// (/meaning there is an incoming connection on the server).
	err = eventLoop.subscribeRead(serverFd)
	if err != nil {
		return nil, err
	}
	return eventLoop, nil
}

// Run runs an event loop. It:
// - runs an event loop in its own goroutine.
// - polls the KQueue for events on the subscribed file descriptors.
// - if the polled event's file descriptor is same as the server's file descriptor: a new client is accepted,
// - else: an existing client for the file descriptor is run.
func (eventLoop *EventLoop) Run() {
	// TODO: Handle client error
	go func() {
		for {
			select {
			case <-eventLoop.stopChannel:
				return
			default:
				events, err := eventLoop.kQueue.Poll(-1)
				if err != nil {
					continue
				}
				for _, event := range events {
					if event.Flags&syscall.EV_EOF == syscall.EV_EOF {
						eventLoop.stopClient(int(event.Ident))
						delete(eventLoop.clients, int(event.Ident))
						continue
					}
					if int(event.Ident) == eventLoop.serverFd {
						if err := eventLoop.acceptClient(); err != nil {
							continue
						}
					} else {
						eventLoop.runClient(int(event.Ident))
					}
				}
			}
		}
	}()
}

// Stop stops the event loop.
func (eventLoop *EventLoop) Stop() {
	close(eventLoop.stopChannel)
	_ = syscall.Close(eventLoop.kQueue.fd)
	for _, client := range eventLoop.clients {
		client.Stop()
	}
}

// subscribeRead subscribes to the given file descriptor using EVFILT_READ filter and an EV_ADD flag which will add the
// file descriptor to the Kernel KQueue when the file descriptor is ready to be read.
func (eventLoop *EventLoop) subscribeRead(fd int) error {
	return eventLoop.kQueue.Subscribe(syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD,
	})
}

// acceptClient accepts a new client (/socket).
// syscall.Accept(..) will not block because the method is called when the non-blocking file descriptor is ready.
func (eventLoop *EventLoop) acceptClient() error {
	fd, _, err := syscall.Accept(eventLoop.serverFd)
	if err != nil {
		return err
	}

	eventLoop.clients[fd] = NewClient(fd, eventLoop.clientHandlers)
	_ = syscall.SetNonblock(fd, true)

	if err := eventLoop.subscribeRead(fd); err != nil {
		return err
	}
	return nil
}

// runClient runs the client for the file descriptor.
func (eventLoop *EventLoop) runClient(fd int) {
	client := eventLoop.clients[fd]
	if client == nil {
		return
	}
	client.Run()
}

// stopClient stops the client corresponding to the file descriptor and closes the descriptor.
func (eventLoop *EventLoop) stopClient(fd int) {
	client := eventLoop.clients[fd]
	if client == nil {
		return
	}
	client.Stop()
	_ = syscall.Close(fd)
}
