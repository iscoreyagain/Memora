//go:build darwin || freebsd || netbsd || openbsd || dragonfly

package io_multiplexing

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/internal/config"
	"golang.org/x/sys/unix"
)

type KQueue struct {
	fd            int
	kqEvents      []unix.Kevent_t
	genericEvents []Event
}

func CreatePoller() (*KQueue, error) {
	epollFD, err := unix.Kqueue()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &KQueue{
		fd:            epollFD,
		kqEvents:      make([]unix.Kevent_t, config.MAXIMUM_CONNECTION),
		genericEvents: make([]Event, config.MAXIMUM_CONNECTION),
	}, nil
}

func (kq *KQueue) Monitor(event Event) error {
	kqEvent := event.toNative(unix.EV_ADD)

	_, err := unix.Kevent(kq.fd, []unix.Kevent_t{kqEvent}, nil, nil)

	return err
}
func (kq *KQueue) Wait() ([]Event, error) {
	n, err := unix.Kevent(kq.fd, nil, kq.kqEvents, nil)
	if err != nil {
		return nil, err
	}

	for i := range n {
		kq.genericEvents[i] = createEvent(kq.kqEvents[i])
	}

	return kq.genericEvents[:n], nil
}

func (kq *KQueue) Close() error {
	return unix.Close(kq.fd)
}
