package single_thread_event_loop

import (
	"fmt"
	"syscall"
	"time"
)

// KQueue represents KQueue for BSD systems.
// fd represents the file descriptor for the kernel KQueue.
type KQueue struct {
	fd       int
	kQEvents []syscall.Kevent_t
}

// NewKQueue creates a new instance of KQueue.
// syscall.Kqueue() creates the Kernel KQueue which can be polled using syscall.Kevent().
func NewKQueue(maxClients int) (*KQueue, error) {
	fd, err := syscall.Kqueue()
	if err != nil {
		return nil, err
	}
	return &KQueue{
		fd:       fd,
		kQEvents: make([]syscall.Kevent_t, maxClients),
	}, nil
}

// Subscribe subscribes to an event of type syscall.Kevent_t.
func (kq *KQueue) Subscribe(event syscall.Kevent_t) error {
	if subscribed, err := syscall.Kevent(
		kq.fd,
		[]syscall.Kevent_t{event},
		nil,
		nil,
	); err != nil || subscribed == -1 {
		return fmt.Errorf("error in subscribing to KQueue: %w", err)
	}
	return nil
}

// Poll polls the Kernel KQueue for the specified duration using Kevent syscall.
// The method blocks until at least one event is triggered or the timeout is reached.
func (kq *KQueue) Poll(timeout time.Duration) ([]syscall.Kevent_t, error) {
	n, err := syscall.Kevent(kq.fd, nil, kq.kQEvents, toTimeSpec(timeout))
	if err != nil {
		return nil, fmt.Errorf("error in KQueue poll: %w", err)
	}
	return kq.kQEvents[:n], nil
}

// Close closes the KQueue.
func (kq *KQueue) Close() error {
	return syscall.Close(kq.fd)
}

// toTimeSpec converts the duration to syscall.Timespec.
func toTimeSpec(duration time.Duration) *syscall.Timespec {
	if duration < 0 {
		return nil
	}
	return &syscall.Timespec{
		Nsec: int64(duration),
	}
}
